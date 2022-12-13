package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/megaease/easeagent-sdk-go/agent"
	"github.com/megaease/easeagent-sdk-go/plugins"
	"github.com/megaease/easeagent-sdk-go/plugins/zipkin"
)

const (
	localHostPort = ":8090"
)

// new tracing agent from yaml file and set host and port of Span.localEndpoint
// By default, use yamlFile="" is Console Reporter for tracing.
// By default, use localHostPort="" is not set host and port of Span.localEndpoint.
var easeagent, _ = agent.NewWithOptions(agent.WithZipkinYAML(os.Getenv("EASEAGENT_CONFIG"), localHostPort))
var tracing = easeagent.GetPlugin(zipkin.NAME).(zipkin.Tracing)

func otherFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("other_function called with method: %s\n", r.Method)
		time.Sleep(50 * time.Millisecond)
	}
}

// http server /some_function span_1
//  - redis get key
// 	- http client WrapHttpRequest(span_1(tracing info))  -> span_2
// 		- http server /other_function 2 -> span_3

func someFunc(url string, client plugins.HTTPDoer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// created span by server middleware
		log.Printf("some_function called with method: %s\n", r.Method)

		// doing some expensive calculations....
		time.Sleep(25 * time.Millisecond)

		log.Printf(url + "/other_function")
		newRequest, err := http.NewRequest("GET", url+"/other_function", nil)
		if err != nil {
			log.Printf("unable to create client: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}

		//send redis span
		redisSpan, _ := tracing.StartMWSpanFromCtx(r.Context(), "redis-get_key", zipkin.Redis)
		if endpoint, err := zipkin.NewEndpoint("redis-local_server", "127.0.0.1:8090"); err == nil {
			redisSpan.SetRemoteEndpoint(endpoint)
		}
		redisSpan.Finish()

		// set server span for parent
		newRequest = easeagent.WrapHTTPRequest(r.Context(), newRequest)
		res, err := client.Do(newRequest)
		if err != nil {
			log.Printf("call to other_function returned error: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}

		// tracing
		_ = res.Body.Close()
	}
}

func main() {
	// initialize router
	router := http.NewServeMux()
	router.HandleFunc("/some_function", someFunc("http://"+localHostPort, easeagent.WrapUserClient(&http.Client{})))
	router.HandleFunc("/other_function", otherFunc())
	http.ListenAndServe(localHostPort, easeagent.WrapUserHandler(router))
}
