// package zipkin

// import (
// 	"net/http"

// 	"github.com/openzipkin/zipkin-go"
// )

// var DEFAULT_AGENT *Agent

// type Agent struct {
// 	options      *Options
// 	zipkinPlugin *ZipkinPlugin
// }

// func Default() *Agent {
// 	return DEFAULT_AGENT
// }

// func InitDefault(hostPort string) {
// 	InitGlobalOptions()
// 	tracingSpec := GlobalOptions.BuildTracingSpec()
// 	//TODO: set tags
// 	InitDefault(tracingSpec)
// 	DEFAULT_AGENT = &Agent{
// 		options:      GlobalOptions,
// 		zipkinPlugin: DEFAULT_PLUGIN,
// 	}
// }
// func CloseDefault() error {
// 	return DEFAULT_AGENT.Close()
// }

// func NewAgent(options *Options) *Agent {
// 	//new tracing instance
// 	return &Agent{
// 		options:      options,
// 		zipkinPlugin: NewPlugin(options.BuildTracingSpec()),
// 	}
// }

// func (agent *Agent) Close() error {
// 	return agent.zipkinPlugin.Close()
// }

// func (agent *Agent) Tracing() *ZipkinTracing {
// 	return agent.zipkinPlugin.Tracing()
// }

// func (agent *Agent) Tracer() *zipkin.Tracer {
// 	return agent.zipkinPlugin.Tracer()
// }

// func (z *Agent) WrapHttpServerHeader(fn http.Handler) http.Handler {
// 	return z.zipkinPlugin.WrapHttpServerHandler(fn)
// }

// func (z *Agent) WrapHttpClient(c *http.Client) HttpClient {
// 	return z.zipkinPlugin.WrapHttpClient(c)
// }
