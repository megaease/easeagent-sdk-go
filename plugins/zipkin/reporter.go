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
	SpanSpec  *SpanSpec
	SenderUrl string
	TlsEnable bool
	TlsKey    string
	TlsCert   string
	TlsCaCert string
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

func httpClient(spec *ReporterSpec) *http.Client {
	if spec.TlsEnable {
		tlsConfig, err := newTLSConfig(spec.TlsCert, spec.TlsKey, spec.TlsCaCert)
		if err != nil {
			exitf("error create tls config: %s", err.Error())
		}
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client := &http.Client{Transport: transport}
		return client
	} else {
		return &http.Client{}
	}
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
