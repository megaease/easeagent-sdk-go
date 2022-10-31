package tracing

import (
	"encoding/json"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
)

type SpanSpec struct {
	Service string
}

type SpanJSONSerializer struct {
	SpanSpec *SpanSpec
}

func (s SpanJSONSerializer) Serialize(spans []*model.SpanModel) ([]byte, error) {
	newSpans := make([]*Span, 0)
	for i := 0; i < len(spans); i++ {
		span := &Span{
			ModelSpanModel: ModelSpanModel(*spans[i]),
			Type:           "log-tracing",
			Service:        s.SpanSpec.Service,
		}
		newSpans = append(newSpans, span)
	}

	return json.Marshal(newSpans)
}

// ContentType returns the ContentType needed for this encoding.
func (SpanJSONSerializer) ContentType() string {
	return "application/json"
}

func SpanSerializer(spanSpec *SpanSpec) reporter.SpanSerializer {
	return &SpanJSONSerializer{
		SpanSpec: spanSpec,
	}
}
