/**
 * Copyright 2022 MegaEase
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *     http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ghodss/yaml"

	"github.com/megaease/easeagent-sdk-go/agent"
	"github.com/megaease/easeagent-sdk-go/plugins"
	"github.com/megaease/easeagent-sdk-go/plugins/easemesh"
	zipkinplugin "github.com/megaease/easeagent-sdk-go/plugins/zipkin"
)

var (
	podServicePort = 80
	podEgressPort  = 13002

	zipKinURL = os.Getenv("ZIPKIN_URL")

	// format: zone.domain.service
	fullServiceName = os.Getenv("SERVICE_NAME")

	// format: service
	internalServiceName string

	zone   string
	domain string
)

const (
	orderSerice                     = "order-mesh"
	restaurantService               = "restaurant-mesh"
	restaurantBeijingAndroidService = "restaurant-mesh-beijing-android"
	restaurantAndroidService        = "restaurant-mesh-android"
	deliveryService                 = "delivery-mesh"
	deliveryBeijingService          = "delivery-mesh-beijing"
	deliveryAndroidService          = "delivery-mesh-android"

	timeFormat = "2006-01-02T15:04:05"
)

type (
	serviceHandler struct {
		urlMutex      sync.Mutex
		restaurantURL string
		deliveryURL   string
	}

	// OrderRequest is the request of order.
	OrderRequest struct {
		OrderID string `json:"order_id"`
		Food    string `json:"food"`
	}

	// OrderResponse is the response of order.
	OrderResponse struct {
		OrderID    string              `json:"order_id"`
		Restaurant *RestaurantResponse `json:"restaurant"`

		ServiceTracings []string `json:"service_tracings,omitempty"`
	}

	// RestaurantRequest is the request of restaurant.
	RestaurantRequest struct {
		OrderID string `json:"order_id"`
		Food    string `json:"food"`
	}

	// RestaurantResponse is the response of restaurant.
	RestaurantResponse struct {
		OrderID      string `json:"order_id"`
		Food         string `json:"food"`
		DeliveryTime string `json:"delivery_time"`

		// Android canary fields.
		Coupon string `json:"coupon,omitempty"`

		ServiceTracings []string `json:"service_tracings,omitempty"`
	}

	// DeliveryRequest is the request of delivery.
	DeliveryRequest struct {
		OrderID string `json:"order_id"`
		Item    string `json:"item"`
	}

	// DeliveryResponse is the response of delivery.
	DeliveryResponse struct {
		OrderID      string `json:"order_id"`
		Item         string `json:"item"`
		DeliveryTime string `json:"delivery_time"`

		// Android canary fields.
		Late *bool `json:"late,omitempty"`

		ServiceTracings []string `json:"service_tracings,omitempty"`
	}
)

var (
	globalHostName   string
	globalAgent      *agent.Agent
	globalHTTPClient plugins.HTTPDoer

	zipkinSpec = zipkinplugin.Spec{
		BaseSpec: plugins.BaseSpec{
			NameField: "zipkin",
			KindField: zipkinplugin.Kind,
		},

		EnableTLS: false,

		EnableBasicAuth: false,

		TracingType: "log-tracing",

		EnableTracing: true,
		SampleRate:    1,
		SharedSpans:   true,
		ID128Bit:      false,
	}
)

func setZipkinSpec() {
	const (
		// NOTE: Just show the available endpoints.
		cnURL         = "https://monitor.megaease.cn:32430/report/application-tracing-log"
		cnInternalURL = "https://172.20.1.116:32330/report/application-tracing-log"
		comURL        = "https://monitor.megaease.com:32430/report/application-tracing-log"
	)

	log.Printf("zipkin url: %s", zipKinURL)

	zipkinSpec.LocalHostport = fmt.Sprintf("%s:80", globalHostName)
	zipkinSpec.OutputServerURL = zipKinURL

	if zipkinSpec.Tags == nil {
		zipkinSpec.Tags = make(map[string]string)
	}

	zipkinSpec.Tags["hostname"] = globalHostName

	if strings.Contains(globalHostName, "shadow") {
		zipkinSpec.Tags["label.local"] = "shadow"
		zipkinSpec.Tags["label.remote"] = "shadow"
	}

	if zipKinURL == "" {
		zipkinSpec.EnableTLS = false
		return
	}

	zipkinSpec.EnableTLS = true

	var certFile, keyFile, caCertFile string
	if strings.Contains(zipKinURL, ".com") {
		certFile = "./tls_cert.com.pem"
		keyFile = "./tls_key.com.key"
		caCertFile = "./tls_ca_cert.com.pem"
	} else {
		certFile = "./tls_cert.cn.pem"
		keyFile = "./tls_key.cn.key"
		caCertFile = "./tls_ca_cert.cn.pem"
	}

	log.Printf("cert file: %s keyfile: %s caCertFile: %s", certFile, keyFile, caCertFile)

	certBuff, err := ioutil.ReadFile(certFile)
	if err != nil {
		exitf("read %s failed: %v", certFile, err)
	}
	zipkinSpec.TLSCert = string(certBuff)

	keyBuff, err := ioutil.ReadFile(keyFile)
	if err != nil {
		exitf("read %s failed: %v", keyFile, err)
	}
	zipkinSpec.TLSKey = string(keyBuff)

	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		exitf("read %s failed: %v", caCertFile, err)
	}
	zipkinSpec.TLSCaCert = string(caCert)
}

func initServiceName() {
	names := strings.Split(fullServiceName, ".")
	if len(names) != 3 {
		exitf("invalid service name: %s", fullServiceName)
		return
	}

	zone, domain, internalServiceName = names[0], names[1], names[2]
}

func completeAnotherServiceName(anotherService string) string {
	return fmt.Sprintf("%s.%s.%s", zone, domain, anotherService)
}

func prefligt() {
	globalHostName, _ = os.Hostname()
	log.Printf("hostname: %s", globalHostName)

	initServiceName()
	log.Printf("full service name: %s", fullServiceName)
	log.Printf("internal service name: %s", internalServiceName)

	zipkinSpec.ServiceName = fullServiceName
	setZipkinSpec()

	if fullServiceName == "" {
		exitf("empty serviceName")
	}

	var agentType string
	switch internalServiceName {
	case orderSerice, deliveryService, deliveryBeijingService, deliveryAndroidService:
		agentType = "GoSDK"
	case restaurantService, restaurantBeijingAndroidService, restaurantAndroidService:
		agentType = "EaseAgent"
	default:
		exitf("unsupport service name: %s", internalServiceName)
	}

	agentConfig := &agent.Config{
		Address: ":9900",
		Plugins: []plugins.Spec{
			easemesh.Spec{
				BaseSpec: plugins.BaseSpec{
					KindField: easemesh.Kind,
					NameField: "easemesh",
				},
				AgentType: agentType,
			},
			zipkinSpec,
		},
	}

	var err error
	globalAgent, err = agent.New(agentConfig)
	if err != nil {
		exitf("create sdk agent failed: %v", err)
	}

	globalHTTPClient = globalAgent.WrapUserClient(http.DefaultClient)
}

func chainReqs(serverReq, clientReq *http.Request) *http.Request {
	return globalAgent.WrapHTTPRequest(serverReq.Context(), clientReq)
}

func main() {
	log.Println("preflight...")
	prefligt()

	serviceServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", podServicePort),
		Handler: globalAgent.WrapUserHandler(newServiceHandler()),
	}

	go func() {
		log.Printf("launch service server...")
		err := serviceServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			exitf("%v", err)
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch

	serviceServer.Shutdown(context.TODO())
}

func exitf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func newServiceHandler() *serviceHandler {
	return &serviceHandler{}
}

func (h *serviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%v", r)
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", r)))
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("read body failed: %v", err)))
		return
	}

	log.Printf("receive %s %s %+v %s", r.Method, r.URL.Path, r.Header, body)

	defer r.Body.Close()

	var resp interface{}

	switch internalServiceName {
	case orderSerice:
		resp, err = h.handleOrder(r, body)
	case restaurantService, restaurantBeijingAndroidService, restaurantAndroidService:
		resp, err = h.handleRestaurant(r, body)
	case deliveryService, deliveryBeijingService, deliveryAndroidService:
		resp, err = h.handleDelivery(r, body)
	default:
		panic(fmt.Errorf("BUG: no correct service"))
	}

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	buff, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")

	// NOTE: For human-readable in the first service.
	if internalServiceName == orderSerice {
		buff, err = yaml.JSONToYAML(buff)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/yaml")
	}

	log.Printf("response: %s", buff)

	w.WriteHeader(200)
	w.Write(buff)
}

func (h *serviceHandler) handleOrder(serverReq *http.Request, body []byte) (interface{}, error) {
	req := &OrderRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}

	restaurantReq := &RestaurantRequest{
		OrderID: req.OrderID,
		Food:    req.Food,
	}

	restaurantReqBody, err := json.Marshal(restaurantReq)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %v", err)
	}

	restaurantURL := fmt.Sprintf("http://%s:%d", completeAnotherServiceName(restaurantService), podEgressPort)
	restaurantClientReq, err := http.NewRequest("POST", restaurantURL, bytes.NewReader(restaurantReqBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}
	restaurantClientReq.Header = serverReq.Header.Clone()
	restaurantClientReq.Header.Set("Content-Type", "application/json")
	restaurantClientReq = chainReqs(serverReq, restaurantClientReq)

	restaurantResp, err := globalHTTPClient.Do(restaurantClientReq)
	if err != nil {
		panic(fmt.Errorf("call restaurant service failed: %v", err))
	}

	if restaurantResp.StatusCode != 200 {
		panic(fmt.Errorf("call restaurant %s failed: status code: %d",
			restaurantURL, restaurantResp.StatusCode))
	}

	restaurantRespBody, err := ioutil.ReadAll(restaurantResp.Body)
	if err != nil {
		panic(fmt.Errorf("read restaurant response body failed: %v", err))
	}

	restaurantResponse := &RestaurantResponse{}
	err = json.Unmarshal(restaurantRespBody, restaurantResponse)
	if err != nil {
		panic(fmt.Errorf("unmarshal restaurant response failed: %v", err))
	}

	resp := &OrderResponse{
		OrderID:    req.OrderID,
		Restaurant: restaurantResponse,
	}

	resp.ServiceTracings = append([]string{globalHostName}, resp.Restaurant.ServiceTracings...)
	resp.Restaurant.ServiceTracings = nil

	return resp, nil
}

func (h *serviceHandler) handleRestaurant(serverReq *http.Request, body []byte) (interface{}, error) {
	req := &RestaurantRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}

	deliveryReq := &DeliveryRequest{
		OrderID: req.OrderID,
		Item:    req.Food,
	}

	deliveryReqBody, err := json.Marshal(deliveryReq)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %v", err)
	}

	deliveryURL := fmt.Sprintf("http://%s:%d", completeAnotherServiceName(deliveryService), podEgressPort)
	deliveryClientReq, err := http.NewRequest("POST", deliveryURL, bytes.NewReader(deliveryReqBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}
	deliveryClientReq.Header = serverReq.Header.Clone()
	deliveryClientReq.Header.Set("Content-Type", "application/json")
	deliveryClientReq = chainReqs(serverReq, deliveryClientReq)

	deliveryResp, err := globalHTTPClient.Do(deliveryClientReq)
	if err != nil {
		panic(fmt.Errorf("call delivery %s failed: %v", deliveryURL, err))
	}

	if deliveryResp.StatusCode != 200 {
		log.Printf("call delivery %s failed: status code: %d", deliveryURL, deliveryResp.StatusCode)
	}

	deliveryRespBody, err := ioutil.ReadAll(deliveryResp.Body)
	if err != nil {
		panic(fmt.Errorf("read delivery response body failed: %v", err))
	}

	deliveryResponse := &DeliveryResponse{}
	err = json.Unmarshal(deliveryRespBody, deliveryResponse)
	if err != nil {
		panic(fmt.Errorf("unmarshal delivery response failed: %v", err))
	}

	deliveryTime := deliveryResponse.DeliveryTime
	// beijing-android restaurant service
	if strings.Contains(globalHostName, "beijing") &&
		strings.Contains(globalHostName, "android") {
		deliveryTime += " (cook duration: 5m)"
	}

	resp := &RestaurantResponse{
		OrderID:      req.OrderID,
		Food:         req.Food,
		DeliveryTime: deliveryTime,
	}

	// android restaurant service
	if strings.Contains(globalHostName, "android") &&
		!strings.Contains(globalHostName, "beijing") {
		if deliveryResponse.Late != nil && *deliveryResponse.Late {
			resp.Coupon = "$5"
		}
	}

	resp.ServiceTracings = append([]string{globalHostName}, deliveryResponse.ServiceTracings...)

	return resp, nil
}

func (h *serviceHandler) handleDelivery(serverReq *http.Request, body []byte) (interface{}, error) {
	log.Printf("header: %+v, body: %s", serverReq.Header, body)

	req := &DeliveryRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}

	deliveryTime := time.Now().Add(10 * time.Minute).Local().Format(timeFormat)

	// beijing delivery service
	if strings.Contains(globalHostName, "beijing") &&
		!strings.Contains(globalHostName, "android") {
		deliveryTime += " (road duration: 7m)"
	}

	resp := &DeliveryResponse{
		OrderID:      req.OrderID,
		Item:         req.Item,
		DeliveryTime: deliveryTime,
	}

	if strings.Contains(globalHostName, "android") &&
		!strings.Contains(globalHostName, "beijing") {
		late := true
		resp.Late = &late
	}

	// NOTE: Make tracing more readable
	time.Sleep(10 * time.Millisecond)

	resp.ServiceTracings = append([]string{globalHostName}, resp.ServiceTracings...)

	return resp, nil
}
