package zipkin

import (
	"fmt"

	"github.com/megaease/easeagent-sdk-go/plugins"
)

const (
	// Kind is the kind of Zipkin plugin.
	Kind = "Zipkin"
	NAME = "Zipkin"
)

type (
	// Spec is the Zipkin spec.
	Spec struct {
		plugins.BaseSpec `json:",inline"`

		OutputServerURL string `json:"reporter.output.server"`

		EnableTLS bool   `json:"reporter.output.server.tls.enable"`
		TLSKey    string `json:"reporter.output.server.tls.key"`
		TLSCert   string `json:"reporter.output.server.tls.cert"`
		TLSCaCert string `json:"reporter.output.server.tls.ca_cert"`

		EnableBasicAuth bool   `json:"reporter.output.server.auth.enable"`
		Username        string `json:"reporter.output.server.auth.username"`
		Password        string `json:"reporter.output.server.auth.password"`

		ServiceName string            `json:"service_name"`
		TracingType string            `json:"tracing_type"`
		Hostport    string            `json:"hostport"`
		Tags        map[string]string `json:"tags"`

		EnableTracing bool    `json:"tracing.enable" jsonschema:"required,minimum=0,maximum=1"`
		SampleRate    float64 `json:"tracing.sample.rate" jsonschema:"required,minimum=0,maximum=1"`
		SharedSpans   bool    `json:"tracing.shared.spans"`
		ID128Bit      bool    `json:"tracing.id128bit"`
	}
)

// DefaultSpec returns the default spec of Zipkin.
func DefaultSpec() plugins.Spec {
	return Spec{
		BaseSpec: plugins.BaseSpec{
			KindField: Kind,
			NameField: NAME,
		},
		OutputServerURL: "https://127.0.0.1:8080/report",

		EnableTLS: false,

		EnableBasicAuth: false,

		EnableTracing: true,
		ServiceName:   "default-service",
		TracingType:   "log-tracing",
		Hostport:      "127.0.0.1:80",

		SampleRate:  1,
		SharedSpans: true,
		ID128Bit:    false,
		Username:    "",
		Password:    "",
	}
}

// Validate validates the Zipkin spec.
func (spec Spec) Validate() error {
	if spec.EnableTLS {
		if len(spec.TLSKey) == 0 || len(spec.TLSCert) == 0 || len(spec.TLSCaCert) == 0 {
			return fmt.Errorf("key, cert, cacert are not all specified")
		}
	}

	if spec.EnableBasicAuth {
		if spec.Username == "" || spec.Password == "" {
			return fmt.Errorf("username and password are not all specified")
		}
	}

	return nil
}
