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
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/stretchr/testify/assert"
)

func TestMarshalJSONAndUnmarshalJSON(t *testing.T) {
	spanModel := &model.SpanModel{
		SpanContext: model.SpanContext{
			ID: 111,
		},
		Duration:  time.Second,
		Timestamp: time.Now(),
		Name:      "testName",
	}
	span := &Span{
		SpanModel: SpanModel(*spanModel),
		Service:   "testServiceName",
		Type:      "log-tracing",
	}
	jsonBytes, err := json.Marshal(span)
	assert.Nil(t, err)
	jsonStr := string(jsonBytes)
	assert.True(t, strings.Contains(jsonStr, "\"testServiceName\""))
	assert.True(t, strings.Contains(jsonStr, "\"log-tracing\""))
	assert.True(t, strings.Contains(jsonStr, "\"timestamp\""))
	assert.True(t, strings.Contains(jsonStr, "\"duration\""))
	fmt.Println(jsonStr)

	var result Span
	err = json.Unmarshal(jsonBytes, &result)
	assert.Nil(t, err)
	assert.Equal(t, span.Service, result.Service)
	assert.Equal(t, span.Type, result.Type)
	assert.Equal(t, span.Timestamp.Round(time.Microsecond).UnixNano(), result.Timestamp.Round(time.Microsecond).UnixNano())
	assert.Equal(t, span.Duration, result.Duration)
	assert.Equal(t, strings.ToLower(span.Name), result.Name)

}
