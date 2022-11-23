package zipkin

import (
	"encoding/json"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
)

type spanJSONSerializer struct {
	serviceName string
	tracingType string
}

func (s spanJSONSerializer) Serialize(spans []*model.SpanModel) ([]byte, error) {
	newSpans := make([]*Span, 0)
	for i := 0; i < len(spans); i++ {
		span := s.WarpSpan(spans[i])
		newSpans = append(newSpans, span)
	}
	return json.Marshal(newSpans)
}

func (s spanJSONSerializer) WarpSpan(span *model.SpanModel) *Span {
	return &Span{
		SpanModel: s.getSpanModel(span),
		Type:      s.tracingType,
		Service:   s.serviceName,
	}
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

func newSpanSerializer(spec Spec) *spanJSONSerializer {
	return &spanJSONSerializer{
		serviceName: spec.ServiceName,
		tracingType: spec.TracingType,
	}
}
