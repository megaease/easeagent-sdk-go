# easeagent-sdk-go

A lightweight & opening Go SDK for Cloud-Native and APM system
## Overview

- EaseAgent SDK can collect distributed application tracing, which could be used in the APM system and improve the observability of a distributed system. for the tracing, EaseAgent SDK follows the [Google Dapper](https://research.google/pubs/pub36356/) paper. 
- EaseAgent SDK also can work with Cloud-Native architecture. For example, it can help Service Mesh (especially for [EaseMesh](https://github.com/megaease/easemesh/) ) to do some control panel work.
- EaseAgent SDK also can work with [MegaCloud](https://cloud.megaease.com/). For example, it can monitor for service by Go Docker APP.

### Principles
- Safe to Go application/service.
- Lightweight and very low CPU, memory, and I/O resource usage.
- Highly extensible, users can easily do extensions through the api
- Design for Micro-Service architecture, collecting the data from a service perspective.

## Features
* Easy to use. It is right out of the box for Http Server Tracing.
  * Collecting Tracing Logs.
    * Http Server
    * Http Client
    * Supplying the `health check` endpoint
  * Decorate the Span API for Middleware

* Data Reports
  * Console Reporter.
  * Http Reporter.

* Standardization
    * The tracing data format is fully compatible with the Zipkin data format.

## QuickStart
### Init Agent
#### 1. Get SDK
```shell
go get github.com/megaease/easeagent-sdk-go
```
#### 2. Import package

```go
import (
	"github.com/megaease/easeagent-sdk-go/agent"
    "github.com/megaease/easeagent-sdk-go/plugins/zipkin"
)
```
#### 3. New Agent
```go
var easeagent = newAgent(hostPort)

// new agent
func newAgent(hostport string) *agent.Agent {
    zipkinSpec = zipkin.DefaultSpec().(zipkin.Spec)
	zipkinSpec.OutputServerURL = "" // report to log when output server is ""
	zipkinSpec.Hostport = hostport
	agent, err := agent.New(&agent.Config{
		Plugins: []plugins.Spec{
			zipkinSpec,
		},
	})
	if err != nil {
		fmt.Fprintf("new agent fail: %v", err)
	    os.Exit(1)
	}
	return agent
}
```
### Warp Http Server
```go
func otherFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("other_function called with method: %s\n", r.Method)
		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	// initialize router
	router := http.NewServeMux()
	router.HandleFunc("/other_function", otherFunc())
	http.ListenAndServe(hostPort, easeagent.WrapUserHandler(router))
}
```
## Example

1. [http example](./example/http/main.go)

2. [mesh example](./example/mesh/main.go)

## About MegaCloud 
1. [Use SDK in MegaCloud](./doc/how-to-use.md)
2. Get MegaCloud Config. [About MegaCloud Config](./doc/megacloud-config.md)
3. [Decorate the Span](./doc/middleware-span.md). please use api: `zipkin.Tracing.StartMWSpan` and `zipkin.Tracing.StartMWSpanFromCtx` for Decorate Span.


