package zipkin

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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

func TestYamlToSpac(t *testing.T) {
	yamlContext := getYaml()
	var body map[string]interface{}
	err := yaml.Unmarshal([]byte(yamlContext), &body)
	assert.Nil(t, err)

	bodyJSON, err := json.Marshal(body)
	assert.Nil(t, err)
	var spec Spec
	err = json.Unmarshal(bodyJSON, &spec)
	assert.Nil(t, err)
	assert.Equal(t, "demo.go.test-service", spec.ServiceName)
	assert.Equal(t, "log-tracing", spec.TracingType)
	assert.True(t, spec.EnableTracing)
	assert.Equal(t, 0.5, spec.SampleRate)
	assert.True(t, spec.SharedSpans)
	assert.False(t, spec.ID128Bit)
	assert.Equal(t, "http://localhost:9411/api/v2/spans", spec.OutputServerURL)
	assert.True(t, spec.EnableTLS)
	assert.Equal(t, "----------- key -----------\nkey\n----------- key end -----------\n", spec.TLSKey)
	assert.Equal(t, "----------- cert -----------\ncert\n----------- cert end -----------\n", spec.TLSCert)
	assert.Equal(t, "----------- ca_cert -----------\nca_cert\n----------- ca_cert end -----------\n", spec.TLSCaCert)
	assert.False(t, spec.EnableBasicAuth)
	assert.Equal(t, "test_user", spec.Username)
	assert.Equal(t, "test_password", spec.Password)
}

func getYaml() string {
	return `service_name: demo.go.test-service
tracing_type: log-tracing
tracing.enable: true
tracing.sample.rate: 0.5
tracing.shared.spans: true
tracing.id128bit: false
reporter.output.server: http://localhost:9411/api/v2/spans
reporter.output.server.tls.enable: true
reporter.output.server.tls.key: |
  ----------- key -----------
  key
  ----------- key end -----------
reporter.output.server.tls.cert: |
  ----------- cert -----------
  cert
  ----------- cert end -----------
reporter.output.server.tls.ca_cert: |
  ----------- ca_cert -----------
  ca_cert
  ----------- ca_cert end -----------
reporter.output.server.auth.enable: false
reporter.output.server.auth.username: test_user
reporter.output.server.auth.password: test_password
`
}
