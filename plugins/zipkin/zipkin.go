package zipkin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/megaease/easemesh/easeagent-sdk-go/plugins"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

// DefaultSpec returns the default spec of EaseMesh.
func DefaultSpec() plugins.Spec {
	return LoadGlobalOptions()
}

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
	// EaseMesh is the EaseMesh dedicated plugin.
	// EaseMesh struct {
	// 	spec Spec

	// 	agentInfo []byte
	// 	headers   atomic.Value // type: []string
	// }

	// // AgentInfo stores agent information.
	// AgentInfo struct {
	// 	Type    string `json:"type"`
	// 	Version string `json:"version"`
	// }

	// // AgentConfig is the config pushed to agent.
	// AgentConfig struct {
	// 	Headers string `json:"easeagent.progress.forwarded.headers"`
	// }
	HandlerWrapper struct {
		handlerFunc http.HandlerFunc
	}
)

func New(spec plugins.Spec) (plugins.Plugin, error) {
	if option, ok := spec.(*Options); ok {
		zipkinPlugin := NewPlugin(option.BuildTracingSpec())
		return zipkinPlugin, nil
	}
	return nil, fmt.Errorf("spec must be *zipkin.Options")
}

var DEFAULT_PLUGIN *ZipkinPlugin

type ZipkinPlugin struct {
	spec    *TracingSpec
	tracing *ZipkinTracing
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
func (h *HandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handlerFunc(w, r)
}

func (z *ZipkinPlugin) WrapUserHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	hander := zipkinhttp.NewServerMiddleware(
		z.tracing.tracer, zipkinhttp.TagResponseSize(true),
	)
	return hander(&HandlerWrapper{
		handlerFunc: handlerFunc,
	}).ServeHTTP
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
