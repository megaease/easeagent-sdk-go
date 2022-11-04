package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/megaease/easeagent-sdk-go/agent"
	"github.com/megaease/easeagent-sdk-go/plugins"
	"github.com/megaease/easeagent-sdk-go/plugins/zipkin"
)

const (
	hostPort = ":8090"
)

var easeagent = newAgent(hostPort)

// new agent
func newAgent(hostport string) *agent.Agent {
	// fileConfigPath := os.Getenv("MEGAEASE_SDK_CONFIG_FILE")
	// if fileConfigPath == "" {
	// 	fileConfigPath = "/megaease/sdk/agent.yml"
	// }
	// zipkinSpec, err := zipkin.LoadSpecFromYamlFile(fileConfigPath)
	// exitfIfErr(err, "new zipkin spec fail: %v", err)
	zipkinSpec := zipkin.DefaultSpec().(zipkin.Spec)
	zipkinSpec.OutputServerURL = "" // report to log when output server is ""
	zipkinSpec.Hostport = hostport
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

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func otherFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("other_function called with method: %s\n", r.Method)
		time.Sleep(50 * time.Millisecond)
	}
}

// http server /some_function span_1
// 	- http client WrapHttpRequest(span_1(tracing info))  -> span_2
// 		- http server /other_function 2 -> span_3

func someFunc(url string, client plugins.HTTPDoer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// created span by server middleware
		log.Printf("some_function called with method: %s\n", r.Method)

		// doing some expensive calculations....
		time.Sleep(25 * time.Millisecond)

		log.Printf(url + "/other_function")
		newRequest, err := http.NewRequest("GET", url+"/other_function", nil)
		if err != nil {
			log.Printf("unable to create client: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		// span := zipkin.SpanFromContext(current)
		// ctx := zipkin.NewContext(req.Context(), span)
		// newRequest = req.WithContext(ctx)

		// set server span for parent
		newRequest = easeagent.WrapHTTPRequest(r.Context(), newRequest)
		res, err := client.Do(newRequest)
		if err != nil {
			log.Printf("call to other_function returned error: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return

		}
		_ = res.Body.Close()
	}
}

func main() {
	// initialize router
	router := http.NewServeMux()
	router.HandleFunc("/hello", hello)
	router.HandleFunc("/headers", headers)
	router.HandleFunc("/some_function", someFunc("http://"+hostPort, easeagent.WrapUserClient(&http.Client{})))
	router.HandleFunc("/other_function", otherFunc())
	http.ListenAndServe(hostPort, easeagent.WrapUserHandler(router))
}
