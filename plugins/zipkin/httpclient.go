package zipkin

import (
	"context"
	"log"
	"net/http"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

var DEFAULT_HTTP_CLIENT HttpClient

type HttpClient interface {
	Do(current context.Context, request *http.Request) (*http.Response, error)
}

type HttpClientImpl struct {
	*zipkinhttp.Client
}

func NewHttpClient(client *zipkinhttp.Client) HttpClient {
	return &HttpClientImpl{
		client,
	}
}

func (c *HttpClientImpl) Do(current context.Context, request *http.Request) (*http.Response, error) {
	span := zipkin.SpanFromContext(current)
	ctx := zipkin.NewContext(request.Context(), span)

	newRequest := request.WithContext(ctx)

	var res *http.Response
	res, err := c.Client.DoWithAppSpan(newRequest, request.Method+" "+request.URL.Path)
	if err != nil {
		log.Printf("call to other_function returned error: %+v\n", err)
		return nil, err
	}
	return res, nil
}
