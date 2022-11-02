package zipkin

import (
	"context"

	"github.com/openzipkin/zipkin-go"
)

func SpanFromContext(ctx context.Context) zipkin.Span {
	return zipkin.SpanFromContext(ctx)
}
