package zipkin

import (
	"fmt"
	"net"
	"time"

	zipkingo "github.com/openzipkin/zipkin-go"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
)

var DEFAULT_TRACING *ZipkinTracing

type TracingSpec struct {
	HostPort           string
	ServiceName        string
	TracingEnable      bool
	TracingSampleRate  float64
	TracingSharedSpans bool
	TracingID128Bit    bool
	TracingTags        map[string]string

	ReporterSpec *ReporterSpec
}

type ZipkinTracing struct {
	spec     *TracingSpec
	endpoint *model.Endpoint
	reporter reporter.Reporter
	tracer   *zipkin.Tracer
}

func DefaultTracing() *ZipkinTracing {
	return DEFAULT_TRACING
}

func InitDefaultTracing(spec *TracingSpec) {
	DEFAULT_TRACING = NewTracing(spec)
}

func CloseDefaultTracing() error {
	return DEFAULT_TRACING.Close()
}

func NewTracing(spec *TracingSpec) *ZipkinTracing {
	endpoint, err := newEndpoint(spec.ServiceName, spec.HostPort)
	exitfIfErr(err, "error creating zipkin endpoint: %v", err)

	reporter := NewReporter(spec.ReporterSpec)
	sampler, err := zipkingo.NewBoundarySampler(spec.TracingSampleRate, time.Now().Unix())
	exitfIfErr(err, "new sampler error: %v", err)

	tracer, err := zipkin.NewTracer(reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithTags(spec.TracingTags),
		zipkingo.WithSampler(sampler),
		zipkingo.WithSharedSpans(spec.TracingSharedSpans),
		zipkingo.WithTraceID128Bit(spec.TracingID128Bit),
	)
	exitfIfErr(err, "tracing init failed: %v", err)
	zipkinTracing := &ZipkinTracing{
		spec:     spec,
		endpoint: endpoint,
		reporter: reporter,
		tracer:   tracer,
	}
	return zipkinTracing
}

func newEndpoint(serviceName string, hostPort string) (*model.Endpoint, error) {
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return nil, err
	}
	if host != "" {
		return zipkin.NewEndpoint(serviceName, hostPort)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				host = ipnet.IP.String()
				break
			}
		}
	}
	return zipkin.NewEndpoint(serviceName, fmt.Sprintf("%s:%s", host, port))
}

func (t *ZipkinTracing) Close() error {
	return t.reporter.Close()
}

func (t *ZipkinTracing) Tracer() *zipkin.Tracer {
	return t.tracer
}
