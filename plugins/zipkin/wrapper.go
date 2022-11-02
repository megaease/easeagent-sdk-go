package zipkin

import (
	"net/http"

	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

type (
	HttpHandlerWrapper struct {
		handlerFunc http.HandlerFunc
	}

	HttpClientWrapper struct {
		client *zipkinhttp.Client
	}
)

func (h *HttpHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handlerFunc(w, r)
}

func (c *HttpClientWrapper) Do(req *http.Request) (*http.Response, error) {
	return c.client.DoWithAppSpan(req, req.Method+" "+req.URL.Path)
}
