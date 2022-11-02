package zipkin

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/megaease/easeagent-sdk-go/plugins"
)

const (
	// Kind is the kind of Zipkin plugin.
	Kind = "ZipKin"
)

// Spec is the Zipkin spec.
type Spec struct {
	plugins.BaseSpec              `json:",inline"`
	TracingEnable                 bool    `default:"true" json:"enable"`
	TracingSampleRate             float64 `default:"1" json:"tracing.sample.rate" jsonschema:"required,minimum=0,maximum=1"`
	TracingSharedSpans            bool    `default:"true" json:"tracing.shared.spans"`
	TracingID128Bit               bool    `default:"true" json:"tracing.id128bit"`
	ReporterOutputServer          string  `json:"reporter.output.server"`
	ReporterOutputServerTlsEnable bool    `default:"false" json:"reporter.output.server.tls.enable"`
	ReporterOutputServerTlsKey    string  `json:"reporter.output.server.tls.key"`
	ReporterOutputServerTlsCert   string  `json:"reporter.output.server.tls.cert"`
	ReporterOutputServerTlsCaCert string  `json:"reporter.output.server.tls.ca_cert"`
	ReporterTracingSenderUrl      string  `json:"reporter.tracing.sender.url"`
	ReporterAuthEnable            bool    `default:"false" json:"reporter.output.server.auth.enable"`
	ReporterAuthUser              string  `json:"reporter.output.server.auth.user"`
	ReporterAuthPassword          string  `json:"reporter.output.server.auth.password"`
	HostPort                      string
	HomeDir                       string `json:"home-dir" long:"home-dir" description:"Path to the home directory."`
}

func NewSpec() *Spec {
	o := &Spec{
		BaseSpec: plugins.BaseSpec{
			KindField: Kind,
			NameField: "demo.demo.easeagent-sdk-go-service",
		},
		TracingEnable:                 true,
		TracingSampleRate:             1,
		TracingSharedSpans:            true,
		TracingID128Bit:               true,
		ReporterOutputServer:          "https://127.0.0.1:8080/report",
		ReporterOutputServerTlsEnable: false,
		ReporterOutputServerTlsKey:    "-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----",
		ReporterOutputServerTlsCert:   "-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----",
		ReporterOutputServerTlsCaCert: "-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----",
		ReporterTracingSenderUrl:      "/application-tracing-log",
		ReporterAuthEnable:            false,
		ReporterAuthUser:              "",
		ReporterAuthPassword:          "",
		HostPort:                      ":8080",
	}
	var err error
	o.HomeDir, err = filepath.Abs(path.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("failed to identify the full home dir: %v", err)
	}

	return o
}

// DefaultSpec returns the default spec of EaseMesh.
func DefaultSpec() plugins.Spec {
	return NewSpec()
}

// Validate validates the Zipkin spec.
func (spec Spec) Validate() error {
	if spec.Name() == "" {
		return fmt.Errorf("name must not be empty")
	}
	if spec.ReporterOutputServerTlsEnable && (spec.ReporterOutputServerTlsKey == "" || spec.ReporterOutputServerTlsCert == "" || spec.ReporterOutputServerTlsCaCert == "") {
		return fmt.Errorf("tls key,cert,cacert must not be empty when tls enable")
	}
	if spec.ReporterAuthEnable && (spec.ReporterAuthUser == "" || spec.ReporterAuthPassword == "") {
		return fmt.Errorf("auth user and password must not be empty when auth enable")
	}
	return nil
}

func (spec *Spec) BuildReporterSpec() *ReporterSpec {
	return &ReporterSpec{
		SpanSpec: &SpanSpec{
			Service: spec.Name(),
		},
		SenderUrl:    spec.ReporterOutputServer + spec.ReporterTracingSenderUrl,
		TlsEnable:    spec.ReporterOutputServerTlsEnable,
		TlsKey:       spec.ReporterOutputServerTlsKey,
		TlsCert:      spec.ReporterOutputServerTlsCert,
		TlsCaCert:    spec.ReporterOutputServerTlsCaCert,
		AuthEnable:   spec.ReporterAuthEnable,
		AuthUser:     spec.ReporterAuthUser,
		AuthPassword: spec.ReporterAuthPassword,
	}
}

func (spec *Spec) BuildTracingSpec() *TracingSpec {
	return &TracingSpec{
		HostPort:           spec.HostPort,
		ServiceName:        spec.Name(),
		TracingEnable:      spec.TracingEnable,
		TracingSampleRate:  spec.TracingSampleRate,
		TracingSharedSpans: spec.TracingSharedSpans,
		TracingID128Bit:    spec.TracingID128Bit,
		TracingTags:        make(map[string]string),
		ReporterSpec:       spec.BuildReporterSpec(),
	}
}
