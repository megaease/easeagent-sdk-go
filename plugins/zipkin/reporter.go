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

package zipkin

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinHttpReporter "github.com/openzipkin/zipkin-go/reporter/http"
)

type (
	// AuthTransport is a http.RoundTripper that adds basic auth to requests.
	AuthTransport struct {
		username string
		password string

		next http.RoundTripper
	}

	// logReporter will send spans to the default Go Logger.
	logReporter struct {
		logger     *log.Logger
		serializer *spanJSONSerializer
	}
)

func newReporter(spec Spec) (reporter.Reporter, error) {
	if spec.OutputServerURL == "" {
		return newLogReporter(spec), nil
	}

	httpClient, err := newHTTPClient(spec)
	if err != nil {
		return nil, fmt.Errorf("new http client failed: %v", err)
	}

	reporter := zipkinHttpReporter.NewReporter(spec.OutputServerURL,
		zipkinHttpReporter.Client(httpClient),
		zipkinHttpReporter.Serializer(newSpanSerializer(spec)))
	return reporter, nil
}

func newHTTPClient(spec Spec) (*http.Client, error) {
	transport := http.DefaultTransport
	if spec.EnableTLS {
		tlsConfig, err := newTLSConfig([]byte(spec.TLSCert), []byte(spec.TLSKey), []byte(spec.TLSCaCert))
		if err != nil {
			return nil, fmt.Errorf("create tls config failed: %v", err)
		}
		transport = &http.Transport{TLSClientConfig: tlsConfig}
	}
	transport = newAuthTransport(spec, transport)
	return &http.Client{Transport: transport}, nil
}

func newAuthTransport(spec Spec, next http.RoundTripper) http.RoundTripper {
	if !spec.EnableBasicAuth {
		return next
	}

	return &AuthTransport{
		username: spec.Username,
		password: spec.Password,
		next:     next,
	}
}

func newTLSConfig(certPem, keyPem, caCertPem []byte) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, fmt.Errorf("load client cert failed: %v", err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCertPem)
	if !ok {
		return nil, fmt.Errorf("load ca cert failed")
	}

	tlsConfig := tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
	}

	return &tlsConfig, nil
}

// RoundTrip adds basic auth to the request.
func (a *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(a.username, a.password)
	return a.next.RoundTrip(req)
}

// NewReporter returns a new log reporter.
func newLogReporter(spec Spec) reporter.Reporter {
	return &logReporter{
		logger:     log.New(os.Stderr, "", log.LstdFlags),
		serializer: newSpanSerializer(spec),
	}
}

// Send outputs a span to the Go logger.
func (r *logReporter) Send(s model.SpanModel) {
	if b, err := json.MarshalIndent(r.serializer.WarpSpan(&s), "", "  "); err == nil {
		r.logger.Printf("%s:\n%s\n\n", time.Now(), string(b))
	}
}

// Close closes the reporter
func (*logReporter) Close() error { return nil }
