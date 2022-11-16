package zipkin

import (
	"encoding/json"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
)

type spanJSONSerializer struct {
	serviceName string
	tracingType string
}

func (s spanJSONSerializer) Serialize(spans []*model.SpanModel) ([]byte, error) {
	newSpans := make([]*Span, 0)
	for i := 0; i < len(spans); i++ {
		span := &Span{
			SpanModel: s.getSpanModel(spans[i]),
			Type:      s.tracingType,
			Service:   s.serviceName,
		}
		newSpans = append(newSpans, span)
	}
	return json.Marshal(newSpans)
}

func (s spanJSONSerializer) getSpanModel(span *model.SpanModel) SpanModel {
	span.RemoteEndpoint = s.getRemoteEndpoint(span)
	return SpanModel(*span)
}

func (s spanJSONSerializer) getRemoteEndpoint(span *model.SpanModel) *model.Endpoint {
	mwTagValue, ok := span.Tags[MIDDLEWARE_TAG]
	if !ok {
		return span.RemoteEndpoint
	}
	if span.RemoteEndpoint == nil {
		if endpoint, err := zipkin.NewEndpoint(mwTagValue, ""); err == nil {
			return endpoint
		} else {
			return nil
		}
	}
	if span.RemoteEndpoint.ServiceName == "" {
		span.RemoteEndpoint.ServiceName = mwTagValue
	}
	return span.RemoteEndpoint
}

// ContentType returns the ContentType needed for this encoding.
func (s spanJSONSerializer) ContentType() string {
	return "application/json"
}

func newSpanSerializer(spec Spec) reporter.SpanSerializer {
	return &spanJSONSerializer{
		serviceName: spec.ServiceName,
		tracingType: spec.TracingType,
	}
}
