package agent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/openzipkin/zipkin-go"
	zipkingo "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinHttpReporter "github.com/openzipkin/zipkin-go/reporter/http"
	logreporter "github.com/openzipkin/zipkin-go/reporter/log"
)

var DEFAULT_AGENT *Agent

type Agent struct {
	options    *Options
	endpoint   *model.Endpoint
	reporter   reporter.Reporter
	tracer     *zipkin.Tracer
	httpClient *zipkinhttp.Client
}

func Default() *Agent {
	return DEFAULT_AGENT
}

func InitDefault(hostPort string) {
	InitGlobalOptions()
	DEFAULT_AGENT = NewAgent(hostPort, GlobalOptions)
}
func CloseDefault() error {
	return DEFAULT_AGENT.Close()
}

func NewAgent(hostPort string, options *Options) *Agent {
	//new tracing instance
	endpoint, err := zipkin.NewEndpoint(options.Name, hostPort)
	if err != nil {
		exitf("error creating zipkin endpoint: %s", err.Error())
	}

	var reporter reporter.Reporter
	if options.ReporterOutputServer == "" {
		reporter = logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags))
		defer func() {
			_ = reporter.Close()
		}()
	} else {
		traceReporterUrl := options.ReporterOutputServer + options.ReporterTracingSenderUrl
		reporter = zipkinHttpReporter.NewReporter(traceReporterUrl, zipkinHttpReporter.Client(httpClient(options)), zipkinHttpReporter.Serializer(spanSerializer(options)))
		// reporter = zipkinHttpReporter.NewReporter(traceReporterServerURL+traceUrl, zipkinHttpReporter.Client(httpClient()))
	}
	traceTags := make(map[string]string)
	sampler, err := zipkingo.NewBoundarySampler(options.TracingSampleRate, time.Now().Unix())
	if err != nil {
		exitf("new sampler error: %s", err.Error())
	}

	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint), zipkin.WithTags(traceTags), zipkingo.WithSampler(sampler))
	if err != nil {
		exitf("tracing init failed: %s", err.Error())
	}
	agent := &Agent{
		options:  options,
		endpoint: endpoint,
		reporter: reporter,
		tracer:   tracer,
	}
	// agent.prefligt()
	return agent
}

func exitf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func httpClient(options *Options) *http.Client {
	if options.ReporterOutputServerTlsEnable {
		tlsConfig, err := newTLSConfig(options.ReporterOutputServerTlsCert, options.ReporterOutputServerTlsKey, options.ReporterOutputServerTlsCaCert)
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

func spanSerializer(options *Options) reporter.SpanSerializer {
	return &SpanJSONSerializer{
		options: options,
	}
}

type SpanJSONSerializer struct {
	options *Options
}

func (s SpanJSONSerializer) Serialize(spans []*model.SpanModel) ([]byte, error) {
	newSpans := make([]*Span, 0)
	for i := 0; i < len(spans); i++ {
		span := &Span{
			ModelSpanModel: ModelSpanModel(*spans[i]),
			Type:           "log-tracing",
			Service:        s.options.Name,
		}
		newSpans = append(newSpans, span)
	}

	return json.Marshal(newSpans)
}

// ContentType returns the ContentType needed for this encoding.
func (SpanJSONSerializer) ContentType() string {
	return "application/json"
}

func (agent *Agent) Close() error {
	return agent.reporter.Close()
}

func (agent *Agent) Tracer() *zipkin.Tracer {
	return agent.tracer
}

func (agent *Agent) HttpClient() *zipkinhttp.Client {
	// create global zipkin traced http client
	//TODO lock
	if agent.httpClient != nil {
		return agent.httpClient
	}
	client, err := zipkinhttp.NewClient(agent.tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}
	agent.httpClient = client
	return client
}

func (agent *Agent) HttpServerMiddleware() func(http.Handler) http.Handler {
	return zipkinhttp.NewServerMiddleware(
		agent.tracer, zipkinhttp.TagResponseSize(true),
	)
}

func (agent *Agent) HttpRequest(current context.Context, request *http.Request) (*http.Response, error) {
	span := zipkin.SpanFromContext(current)
	ctx := zipkin.NewContext(request.Context(), span)

	newRequest := request.WithContext(ctx)

	var res *http.Response
	res, err := agent.HttpClient().DoWithAppSpan(newRequest, request.Method+" "+request.URL.Path)
	if err != nil {
		log.Printf("call to other_function returned error: %+v\n", err)
		return nil, err
	}
	return res, nil
}
