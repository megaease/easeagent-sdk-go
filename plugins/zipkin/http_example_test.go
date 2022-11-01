package zipkin

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

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

func someFunc(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("some_function called with method: %s\n", r.Method)

		// retrieve span from context (created by server middleware)
		// agent.Default().Tracer().
		span := tracing.SpanFromContext(r.Context())
		span.Tag("custom_key", "some value")

		// doing some expensive calculations....
		time.Sleep(25 * time.Millisecond)
		span.Annotate(time.Now(), "expensive_calc_done")

		log.Printf(url + "/other_function")

		newRequest, err := http.NewRequest("GET", url+"/other_function", nil)
		if err != nil {
			log.Printf("unable to create client: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		res, err := tracing.DEFAULT_HTTP_CLIENT.Do(r.Context(), newRequest)
		if err != nil {
			log.Printf("call to other_function returned error: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return

		}
		_ = res.Body.Close()
	}
}

func Example_main() {
	hostPort := "127.0.0.1:8090" // your host and server port
	agent.InitDefault(hostPort)
	defer agent.CloseDefault()

	// // initialize router
	router := http.NewServeMux()
	router.HandleFunc("/hello", hello)
	router.HandleFunc("/headers", headers)
	router.HandleFunc("/some_function", someFunc("http://"+hostPort))
	router.HandleFunc("/other_function", otherFunc())
	http.ListenAndServe(":8090", agent.Default().WrapHttpServerHeader(router))

	// http.HandleFunc("/hello", hello)
	// http.HandleFunc("/headers", headers)
	// err := http.ListenAndServe(":8090", agent.HttpServerMiddleware()(http.DefaultServeMux))
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
