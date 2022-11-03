package zipkin

import (
	"context"
	"fmt"
	"os"

	"github.com/openzipkin/zipkin-go"
)

func exitfIfErr(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func SpanFromContext(ctx context.Context) zipkin.Span {
	return zipkin.SpanFromContext(ctx)
}
