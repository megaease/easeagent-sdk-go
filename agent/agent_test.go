package agent

import (
	"fmt"
	"testing"
	"time"

	"github.com/openzipkin/zipkin-go/model"
)

func TestSerialize(t *testing.T) {
	serializer := &SpanJSONSerializer{}
	// serializer := reporter.JSONSerializer{}
	spans := make([]*model.SpanModel, 1)
	spans[0] = &model.SpanModel{
		SpanContext: model.SpanContext{},
		Timestamp:   time.Now(),
		Name:        "testName",
	}
	d, _ := serializer.Serialize(spans)
	fmt.Println(string(d))
}
