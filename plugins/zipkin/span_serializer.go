/**
 * Copyright 2022 MegaEase
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *     http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package zipkin

import (
	"encoding/json"

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
	mwTagValue, ok := span.Tags[MiddlewareTag]
	if !ok {
		return span.RemoteEndpoint
	}
	if span.RemoteEndpoint == nil {
		return NewEndpointByName(mwTagValue)
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
