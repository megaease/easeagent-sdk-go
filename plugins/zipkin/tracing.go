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

func NewTracing(spec *TracingSpec) (*ZipkinTracing, error) {
	endpoint, err := newEndpoint(spec.ServiceName, spec.HostPort)
	if err != nil {
		return nil, fmt.Errorf("error creating zipkin endpoint: %v", err)
	}

	reporter, err := NewReporter(spec.ReporterSpec)
	if err != nil {
		return nil, err
	}

	sampler, err := zipkingo.NewBoundarySampler(spec.TracingSampleRate, time.Now().Unix())
	if err != nil {
		return nil, fmt.Errorf("new sampler error: %v", err)
	}

	tracer, err := zipkin.NewTracer(reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithTags(spec.TracingTags),
		zipkingo.WithSampler(sampler),
		zipkingo.WithSharedSpans(spec.TracingSharedSpans),
		zipkingo.WithTraceID128Bit(spec.TracingID128Bit),
	)
	if err != nil {
		return nil, fmt.Errorf("tracing init failed: %v", err)
	}

	zipkinTracing := &ZipkinTracing{
		spec:     spec,
		endpoint: endpoint,
		reporter: reporter,
		tracer:   tracer,
	}
	return zipkinTracing, nil
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
