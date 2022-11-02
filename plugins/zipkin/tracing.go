package zipkin

import (
	"fmt"
	"os"
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
	endpoint, err := zipkin.NewEndpoint(spec.ServiceName, spec.HostPort)
	if err != nil {
		exitf("error creating zipkin endpoint: %s", err.Error())
	}

	reporter := NewReporter(spec.ReporterSpec)
	sampler, err := zipkingo.NewBoundarySampler(spec.TracingSampleRate, time.Now().Unix())
	if err != nil {
		exitf("new sampler error: %s", err.Error())
	}

	tracer, err := zipkin.NewTracer(reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithTags(spec.TracingTags),
		zipkingo.WithSampler(sampler),
		zipkingo.WithSharedSpans(spec.TracingSharedSpans),
		zipkingo.WithTraceID128Bit(spec.TracingID128Bit),
	)
	if err != nil {
		exitf("tracing init failed: %s", err.Error())
	}
	zipkinTracing := &ZipkinTracing{
		spec:     spec,
		endpoint: endpoint,
		reporter: reporter,
		tracer:   tracer,
	}
	return zipkinTracing
}

func exitf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func (t *ZipkinTracing) Close() error {
	return t.reporter.Close()
}

func (t *ZipkinTracing) Tracer() *zipkin.Tracer {
	return t.tracer
}
