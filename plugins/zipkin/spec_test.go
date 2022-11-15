package zipkin

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpac(t *testing.T) {
	bodyJSON := `{"reporter.output.server":"https://127.0.0.1:32430/report/application-tracing-log"}`
	var spec Spec
	err := json.Unmarshal([]byte(bodyJSON), &spec)
	assert.Nil(t, err)
	// assert.isNil(t, "err", err)

	bodyJSON = `{"reporter.output.server":"https://127.0.0.1:32430/report/application-tracing-log", "reporter.output.server.tls.key": "-----BEGIN PRIVATE KEY-----"}`
	var spec2 Spec
	err = json.Unmarshal([]byte(bodyJSON), &spec2)
	assert.Nil(t, err)
}
