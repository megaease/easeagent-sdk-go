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
	hostPort = ":8090" // your server host and port for
)

var easeagent = newAgent(hostPort)
var tracing = easeagent.GetPlugin(zipkin.NAME).(zipkin.Tracing)

// new agent
func newAgent(hostport string) *agent.Agent {
	fileConfigPath := os.Getenv("MEGAEASE_SDK_CONFIG_FILE")
	if fileConfigPath == "" {
		fileConfigPath = "/megaease/sdk/agent.yml"
	}
	spec, err := LoadSpecFromYamlFile(fileConfigPath, hostport)
	zipkinSpec := *spec
	exitfIfErr(err, "new zipkin spec fail: %v", err)
	agent, err := agent.New(&agent.Config{
		Plugins: []plugins.Spec{
			zipkinSpec,
		},
	})
	exitfIfErr(err, "new agent fail: %v", err)
	return agent
}

func exitfIfErr(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func LoadSpecFromYamlFile(filePath string, hostport string) (*zipkin.Spec, error) {
	buff, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read config file :%s failed: %v", filePath, err)
	}
	fmt.Println(string(buff))
	var body map[string]interface{}
	err = yaml.Unmarshal(buff, &body)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml file %s to map failed: %v",
			filePath, err)
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal yaml file %s to json failed: %v",
			filePath, err)
	}
	var spec zipkin.Spec
	err = json.Unmarshal(bodyJson, &spec)
	if err != nil {
		return nil, fmt.Errorf("unmarshal %s to %T failed: %v", bodyJson, spec, err)
	}
	spec.KindField = zipkin.Kind
	spec.NameField = spec.ServiceName
	spec.Hostport = hostport
	return &spec, nil
}
```
### Third: Wrap HTTP

##### 1. Wrap Server Handler 
```go
router := http.NewServeMux()
http.ListenAndServe(hostPort, easeagent.WrapUserHandler(router))
```

##### 2. Wrap Client and Request
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
