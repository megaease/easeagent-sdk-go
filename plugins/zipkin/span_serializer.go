package zipkin

import (
	"encoding/json"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
)

type spanJSONSerializer struct {
	serviceName string
}

func (s spanJSONSerializer) Serialize(spans []*model.SpanModel) ([]byte, error) {
	newSpans := make([]*Span, 0)
	for i := 0; i < len(spans); i++ {
		span := &Span{
			SpanModel: model.SpanModel(*spans[i]),
			Type:      "log-tracing",
			Service:   s.serviceName,
		}
		newSpans = append(newSpans, span)
	}

	return json.Marshal(newSpans)
}

// ContentType returns the ContentType needed for this encoding.
func (s spanJSONSerializer) ContentType() string {
	return "application/json"
}

func newSpanSerializer(spec Spec) reporter.SpanSerializer {
	return &spanJSONSerializer{
		serviceName: spec.ServiceName,
	}
}
