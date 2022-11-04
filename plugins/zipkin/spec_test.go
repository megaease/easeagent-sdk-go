package zipkin

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpac(t *testing.T) {
	bodyJson := `{"reporter.output.server":"https://127.0.0.1:32430/report/application-tracing-log"}`
	var spec Spec
	err := json.Unmarshal([]byte(bodyJson), &spec)
	assert.Nil(t, err)
	// assert.isNil(t, "err", err)

	bodyJson = `{"reporter.output.server":"https://127.0.0.1:32430/report/application-tracing-log", "reporter.output.server.tls.key": "-----BEGIN PRIVATE KEY-----"}`
	var spec2 Spec
	err = json.Unmarshal([]byte(bodyJson), &spec2)
	assert.Nil(t, err)
}
