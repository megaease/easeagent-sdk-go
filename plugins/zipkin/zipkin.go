package zipkin

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/megaease/easemesh/easeagent-sdk-go/plugins"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

func init() {
	cons := &plugins.Constructor{
		Kind:         Kind,
		DefaultSpec:  DefaultSpec,
		SystemPlugin: false,
		NewInstance:  New,
	}
	plugins.Register(cons)
}

type (
	ZipkinPlugin struct {
		spec    *TracingSpec
		tracing *ZipkinTracing
	}
)

func New(spec plugins.Spec) (plugins.Plugin, error) {
	if spec, ok := spec.(*Spec); ok {
		zipkinPlugin := NewPlugin(spec.BuildTracingSpec())
		return zipkinPlugin, nil
	}
	return nil, fmt.Errorf("spec must be *zipkin.Spec")
}

func NewPlugin(spec *TracingSpec) *ZipkinPlugin {
	return &ZipkinPlugin{
		spec:    spec,
		tracing: NewTracing(spec),
	}
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

func (z *ZipkinPlugin) WrapUserHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	hander := zipkinhttp.NewServerMiddleware(
		z.tracing.tracer, zipkinhttp.TagResponseSize(true),
	)
	return hander(&HTTPHandlerWrapper{
		handlerFunc: handlerFunc,
	}).ServeHTTP
}

func (z *ZipkinPlugin) WrapUserClient(c plugins.HTTPDoer) plugins.HTTPDoer {
	if original, ok := c.(*http.Client); ok {
		client, err := zipkinhttp.NewClient(z.tracing.tracer,
			zipkinhttp.WithClient(original),
			zipkinhttp.ClientTrace(z.spec.TracingEnable),
		)
		if err != nil {
			log.Fatalf("unable to create client: %+v\n", err)
		}
		return &HTTPClientWrapper{
			client: client,
		}
	}
	log.Println("can warp plugins.HTTPDoer for zipkin, it must be a *http.Client")
	return c
}

func (z *ZipkinPlugin) WrapUserClientRequest(current context.Context, req *http.Request) *http.Request {
	span := zipkin.SpanFromContext(current)
	ctx := zipkin.NewContext(req.Context(), span)
	return req.WithContext(ctx)
}
