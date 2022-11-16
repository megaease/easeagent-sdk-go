package zipkin

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/megaease/easeagent-sdk-go/plugins"
	"github.com/openzipkin/zipkin-go"
	zipkingo "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
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
	// Zipkin is the Zipkin dedicated plugin.
	Zipkin struct {
		spec Spec

		reporter reporter.Reporter
		tracer   *zipkin.Tracer
	}
)

// New creates a new Zipkin plugin.
func New(pluginSpec plugins.Spec) (plugins.Plugin, error) {
	spec := pluginSpec.(Spec)

	endpoint, err := NewEndpoint(spec.ServiceName, spec.Hostport)
	if err != nil {
		return nil, fmt.Errorf("new endpoint failed: %v", err)
	}

	reporter, err := newReporter(spec)
	if err != nil {
		return nil, fmt.Errorf("new reporter failed: %v", err)
	}

	sampler, err := zipkingo.NewBoundarySampler(spec.SampleRate, time.Now().Unix())
	if err != nil {
		return nil, fmt.Errorf("new sampler failed: %v", err)
	}

	tracer, err := zipkin.NewTracer(reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithTags(spec.Tags),
		zipkingo.WithSampler(sampler),
		zipkingo.WithSharedSpans(spec.SharedSpans),
		zipkingo.WithTraceID128Bit(spec.ID128Bit),
	)
	if err != nil {
		return nil, fmt.Errorf("new tracer failed: %v", err)
	}

	z := &Zipkin{
		spec:     spec,
		tracer:   tracer,
		reporter: reporter,
	}

	return z, nil
}

func NewEndpoint(serviceName string, hostPort string) (*model.Endpoint, error) {
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

func (z *Zipkin) Name() string {
	return z.spec.Name()
}

// Close closes the plugin.
func (z *Zipkin) Close() error {
	return z.reporter.Close()
}

// WrapUserHandlerFunc wraps the user's http handler.
func (z *Zipkin) WrapUserHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	handler := zipkinhttp.NewServerMiddleware(
		z.tracer, zipkinhttp.TagResponseSize(true),
	)
	return handler(&HTTPHandlerWrapper{
		handlerFunc: handlerFunc,
	}).ServeHTTP
}

// WrapUserClient wraps the http client.
func (z *Zipkin) WrapUserClient(c plugins.HTTPDoer) plugins.HTTPDoer {
	if original, ok := c.(*http.Client); ok {
		client, err := zipkinhttp.NewClient(z.tracer,
			zipkinhttp.WithClient(original),
			zipkinhttp.ClientTrace(true),
		)
		if err != nil {
			log.Printf("unable to create client: %+v\n", err)
			return c
		}
		return &HTTPClientWrapper{
			client: client,
		}
	}
	log.Println("can warp plugins.HTTPDoer for zipkin, it must be a *http.Client")
	return c
}

// WrapUserClientRequest wraps the user's http request.
func (z *Zipkin) WrapUserClientRequest(current context.Context, req *http.Request) *http.Request {
	span := zipkin.SpanFromContext(current)
	ctx := zipkin.NewContext(req.Context(), span)
	return req.WithContext(ctx)
}

//start a Span from parent
func (z *Zipkin) StartSpan(parent zipkin.Span, name string, options ...zipkin.SpanOption) zipkin.Span {
	if parent == nil {
		return z.tracer.StartSpan(name, options...)
	}
	options = append(options, zipkin.Parent(parent.Context()))
	return z.tracer.StartSpan(name, options...)
}

//start a Span from context.Context
func (z *Zipkin) StartSpanFromCtx(parent context.Context, name string, options ...zipkin.SpanOption) (zipkin.Span, context.Context) {
	return z.tracer.StartSpanFromContext(parent, name, options...)
}

//start a middleware span from parent
func (z *Zipkin) StartMWSpan(parent zipkin.Span, name string, mwType MiddlewareType, options ...zipkin.SpanOption) zipkin.Span {
	os := make([]zipkin.SpanOption, 0)
	os = append(os, zipkin.Kind(model.Client))
	os = append(os, options...)
	span := z.StartSpan(parent, name, os...)
	span.Tag(MIDDLEWARE_TAG, mwType.TagValue())
	return span
}

//start a middleware span from context.Context
func (z *Zipkin) StartMWSpanFromCtx(parent context.Context, name string, mwType MiddlewareType, options ...zipkin.SpanOption) (zipkin.Span, context.Context) {
	os := make([]zipkin.SpanOption, 0)
	os = append(os, zipkin.Kind(model.Client))
	os = append(os, options...)
	span, ctx := z.StartSpanFromCtx(parent, name, os...)
	span.Tag(MIDDLEWARE_TAG, mwType.TagValue())
	return span, ctx
}
