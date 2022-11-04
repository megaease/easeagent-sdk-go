package zipkin

import (
	"net/http"

	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

type (
	// HTTPHandlerWrapper is the wrapper of http.Handler.
	HTTPHandlerWrapper struct {
		handlerFunc http.HandlerFunc
	}

	// HTTPClientWrapper is the wrapper of http.Client.
	HTTPClientWrapper struct {
		client *zipkinhttp.Client
	}
)

func (h *HTTPHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handlerFunc(w, r)
}

// Do implements plugins.HTTPDoer.
func (c *HTTPClientWrapper) Do(req *http.Request) (*http.Response, error) {
	return c.client.DoWithAppSpan(req, req.Method+" "+req.URL.Path)
}
