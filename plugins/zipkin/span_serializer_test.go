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
	"testing"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	serializer := &spanJSONSerializer{
		serviceName: "testServiceName",
		tracingType: "log-tracing",
	}
	// serializer := reporter.JSONSerializer{}
	spans := make([]*model.SpanModel, 1)
	spans[0] = &model.SpanModel{
		SpanContext: model.SpanContext{},
		Timestamp:   time.Now(),
		Duration:    time.Second,
		Name:        "testName",
	}
	d, _ := serializer.Serialize(spans)
	fmt.Println(string(d))
	var data []map[string]string
	json.Unmarshal(d, &data)
	assert.Equal(t, 1, len(data))
	spanMap := data[0]
	assert.Equal(t, "testServiceName", spanMap["service"])
	assert.Equal(t, "log-tracing", spanMap["type"])
	_, ok := spanMap["timestamp"]
	assert.True(t, ok)
	_, ok = spanMap["duration"]
	assert.True(t, ok)
}
