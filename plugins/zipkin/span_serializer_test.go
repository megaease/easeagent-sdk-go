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
		Name:        "testName",
	}
	d, _ := serializer.Serialize(spans)
	var data []map[string]string
	json.Unmarshal(d, &data)
	assert.Equal(t, 1, len(data))
	spanMap := data[0]
	assert.Equal(t, "testServiceName", spanMap["service"])
	assert.Equal(t, "log-tracing", spanMap["type"])
	assert.NotNil(t, spanMap["timestamp"])
	assert.NotNil(t, spanMap["duration"])
	fmt.Println(string(d))
}
