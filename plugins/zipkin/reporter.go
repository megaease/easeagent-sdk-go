package zipkin

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/openzipkin/zipkin-go/reporter"
	zipkinHttpReporter "github.com/openzipkin/zipkin-go/reporter/http"
	logreporter "github.com/openzipkin/zipkin-go/reporter/log"
)

type ReporterSpec struct {
	SpanSpec     *SpanSpec
	SenderUrl    string
	TlsEnable    bool
	TlsKey       string
	TlsCert      string
	TlsCaCert    string
	AuthEnable   bool
	AuthUser     string
	AuthPassword string
}

type AuthTransport struct {
	user     string
	password string
	next     http.RoundTripper
}

func NewReporter(spec *ReporterSpec) reporter.Reporter {
	var reporter reporter.Reporter
	if spec.SenderUrl == "" {
		reporter = logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags))
		defer func() {
			_ = reporter.Close()
		}()
	} else {
		reporter = zipkinHttpReporter.NewReporter(spec.SenderUrl, zipkinHttpReporter.Client(httpClient(spec)), zipkinHttpReporter.Serializer(SpanSerializer(spec.SpanSpec)))
	}
	return reporter
}

func basicAuthTransport(spec *ReporterSpec, next http.RoundTripper) http.RoundTripper {
	if !spec.AuthEnable {
		return next
	}
	return &AuthTransport{
		user:     spec.AuthUser,
		password: spec.AuthPassword,
		next:     next,
	}
}

func httpClient(spec *ReporterSpec) *http.Client {
	transport := http.DefaultTransport
	if spec.TlsEnable {
		tlsConfig, err := newTLSConfig(spec.TlsCert, spec.TlsKey, spec.TlsCaCert)
		if err != nil {
			exitf("error create tls config: %s", err.Error())
		}
		transport = &http.Transport{TLSClientConfig: tlsConfig}
	}
	transport = basicAuthTransport(spec, transport)
	return &http.Client{Transport: transport}
}

func newTLSConfig(clientCert, clientKey, caCert string) (*tls.Config, error) {
	tlsConfig := tls.Config{InsecureSkipVerify: true}

	// Load client cert
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(caCert))
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err

}

func (a *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(a.user, a.password)
	return a.next.RoundTrip(req)
}
