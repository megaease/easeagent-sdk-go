package tracing

import (
	"log"
	"net/http"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

var DEFAULT_PLUGIN *ZipkinPlugin

type ZipkinPlugin struct {
	spec    *TracingSpec
	tracing *ZipkinTracing
}

func Default() *ZipkinPlugin {
	return DEFAULT_PLUGIN
}

func InitDefault(spec *TracingSpec) {
	InitDefaultTracing(spec)
	DEFAULT_PLUGIN = &ZipkinPlugin{
		tracing: DEFAULT_TRACING,
	}
	DEFAULT_HTTP_CLIENT = DEFAULT_PLUGIN.WrapHttpClient(nil)
}

func NewPlugin(spec *TracingSpec) *ZipkinPlugin {
	return &ZipkinPlugin{
		tracing: NewTracing(spec),
	}
}

func CloseDefault() error {
	return DEFAULT_PLUGIN.Close()
}

func (z *ZipkinPlugin) Tracer() *zipkin.Tracer {
	return z.tracing.tracer
}

func (z *ZipkinPlugin) Close() error {
	return z.tracing.Close()
}

func (z *ZipkinPlugin) Tracing() *ZipkinTracing {
	return z.tracing
}

func (z *ZipkinPlugin) WrapHttpServerHandler(fn http.Handler) http.Handler {
	return zipkinhttp.NewServerMiddleware(
		z.tracing.tracer, zipkinhttp.TagResponseSize(true),
	)(fn)
}

func (z *ZipkinPlugin) WrapHttpClient(c *http.Client) HttpClient {
	client, err := zipkinhttp.NewClient(z.tracing.tracer,
		zipkinhttp.WithClient(c),
		zipkinhttp.ClientTrace(z.spec.TracingEnable),
	)
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}
	return NewHttpClient(client)
}
