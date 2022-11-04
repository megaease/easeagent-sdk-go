package zipkin

import (
	"fmt"

	"github.com/megaease/easeagent-sdk-go/plugins"
)

const (
	// Kind is the kind of Zipkin plugin.
	Kind = "Zipkin"
)

type (
	// Spec is the Zipkin spec.
	Spec struct {
		plugins.BaseSpec `json:",inline"`

		OutputServerURL string `json:"reporter.output.server"`

		EnableTLS bool   `json:"reporter.output.server.tls.enable"`
		TLSKey    []byte `json:"reporter.output.server.tls.key"`
		TLSCert   []byte `json:"reporter.output.server.tls.cert"`
		TLSCaCert []byte `json:"reporter.output.server.tls.ca_cert"`

		EnableBasicAuth bool   `json:"reporter.output.server.auth.enable"`
		Username        string `json:"reporter.output.server.auth.username"`
		Password        string `json:"reporter.output.server.auth.password"`

		ServiceName string            `json:"service_name"`
		Hostport    string            `json:"hostport"`
		Tags        map[string]string `json:"tags"`

		SampleRate  float64 `json:"tracing.sample.rate" jsonschema:"required,minimum=0,maximum=1"`
		SharedSpans bool    `json:"tracing.shared.spans"`
		ID128Bit    bool    `json:"tracing.id128bit"`
	}
)

// DefaultSpec returns the default spec of Zipkin.
func DefaultSpec() plugins.Spec {
	return Spec{
		BaseSpec: plugins.BaseSpec{
			KindField: Kind,
			NameField: "demo.demo.easeagent-sdk-go-service",
		},
		OutputServerURL: "https://127.0.0.1:8080/report",

		EnableTLS: false,

		EnableBasicAuth: false,

		ServiceName: "default-service",
		Hostport:    "127.0.0.1:80",

		SampleRate:  1,
		SharedSpans: true,
		ID128Bit:    true,
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

// func (spec Spec) BuildReporterSpec() *ReporterSpec {
// 	reporterSpec := &ReporterSpec{
// 		SpanSpec: &SpanSpec{
// 			Service: spec.Name(),
// 		},
// 		TLSEnable:    spec.EnableTLS,
// 		TLSKey:       spec.TLSKey,
// 		TLSCert:      spec.TLSCert,
// 		TLSCaCert:    spec.TLSCaCert,
// 		AuthEnable:   spec.EnableBasicAuth,
// 		AuthUser:     spec.Username,
// 		AuthPassword: spec.Password,
// 	}
// 	if spec.OutputServerURL != "" {
// 		reporterSpec.SenderURL = spec.OutputServerURL + spec.ReporterTracingSenderURL
// 	} else {
// 		reporterSpec.SenderURL = ""
// 	}
// 	return reporterSpec
// }

// func (spec *Spec) BuildTracingSpec() *TracingSpec {
// 	return &TracingSpec{
// 		HostPort:           spec.Hostport,
// 		ServiceName:        spec.Name(),
// 		TracingEnable:      spec.TracingEnable,
// 		TracingSampleRate:  spec.SampleRate,
// 		TracingSharedSpans: spec.SharedSpans,
// 		TracingID128Bit:    spec.ID128Bit,
// 		TracingTags:        make(map[string]string),
// 		ReporterSpec:       spec.BuildReporterSpec(),
// 	}
// }
