### First: Get SDK
```bash
cd {project}
go get github.com/megaease/easeagent-sdk-go
```

### Second: Init Agent

##### 1. Import package
    
```go
import (
	"github.com/megaease/easeagent-sdk-go/agent"
	"github.com/megaease/easeagent-sdk-go/plugins"
	"github.com/megaease/easeagent-sdk-go/plugins/zipkin"
	"gopkg.in/yaml.v2"
)
```

##### 2. Init Agent
You can load spec then new Agent like below code:
```go
const (
	localHostPort = ":8090" // your server host and port for
)
// new tracing agent from yaml file and set host and port of Span.localEndpoint
// By default, use yamlFile="" is use easemesh.DefaultSpec() and Console Reporter for tracing.
// By default, use localHostPort="" is not set host and port of Span.localEndpoint.
var easeagent, _ = agent.NewWithOptions(agent.WithYAML(os.Getenv("EASEAGENT_CONFIG"), localHostPort))
var tracing = easeagent.GetPlugin(zipkin.NAME).(zipkin.Tracing)
```
### Third: Wrapping HTTP

##### 1. Wrapping Server Handler 
```go
router := http.NewServeMux()
http.ListenAndServe(hostPort, easeagent.WrapUserHandler(router))
```

##### 2. Wrapping Client and Request
```go
client := easeagent.WrapUserClient(&http.Client{})
newRequest, err := http.NewRequest("GET", url+"/other_function", nil)
newRequest = easeagent.WrapHTTPRequest(serverRequest.Context(), newRequest)
res, err := client.Do(newRequest)
```

##### 3. Tracing middleware
```go
//send redis span
redisSpan, _ := tracing.StartMWSpanFromCtx(r.Context(), "redis-get_key", zipkin.Redis)
if endpoint, err := zipkin.NewEndpoint("redis-local_server", "127.0.0.1:8090"); err == nil {
    redisSpan.SetRemoteEndpoint(endpoint)
}
redisSpan.Finish()
```
## Example
[HTTP example](../example/http/main.go)
