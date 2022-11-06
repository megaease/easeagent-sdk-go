package zipkin

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/openzipkin/zipkin-go/reporter"
	zipkinHttpReporter "github.com/openzipkin/zipkin-go/reporter/http"
	logreporter "github.com/openzipkin/zipkin-go/reporter/log"
)

type (
	// AuthTransport is a http.RoundTripper that adds basic auth to requests.
	AuthTransport struct {
		username string
		password string

		next http.RoundTripper
	}
)

func newReporter(spec Spec) (reporter.Reporter, error) {
	if spec.OutputServerURL == "" {
		return logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags)), nil
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
