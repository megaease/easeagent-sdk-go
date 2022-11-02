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
	ReporterOutputServerTlsEnable bool    `json:"reporter.output.server.tls.enable"`
	ReporterOutputServerTlsKey    string  `json:"reporter.output.server.tls.key"`
	ReporterOutputServerTlsCert   string  `json:"reporter.output.server.tls.cert"`
	ReporterOutputServerTlsCaCert string  `json:"reporter.output.server.tls.ca_cert"`
	ReporterTracingSenderUrl      string  `json:"reporter.tracing.sender.url"`
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
func (opt Spec) Validate() error {
	if opt.Name() == "" {
		return fmt.Errorf("name must not be empty")
	}
	return nil
}

func (opt *Spec) BuildReporterSpec() *ReporterSpec {
	return &ReporterSpec{
		SpanSpec: &SpanSpec{
			Service: opt.Name(),
		},
		SenderUrl: opt.ReporterOutputServer + opt.ReporterTracingSenderUrl,
		TlsEnable: opt.ReporterOutputServerTlsEnable,
		TlsKey:    opt.ReporterOutputServerTlsKey,
		TlsCert:   opt.ReporterOutputServerTlsCert,
		TlsCaCert: opt.ReporterOutputServerTlsCaCert,
	}
}

func (opt *Spec) BuildTracingSpec() *TracingSpec {
	return &TracingSpec{
		HostPort:           opt.HostPort,
		ServiceName:        opt.Name(),
		TracingEnable:      opt.TracingEnable,
		TracingSampleRate:  opt.TracingSampleRate,
		TracingSharedSpans: opt.TracingSharedSpans,
		TracingID128Bit:    opt.TracingID128Bit,
		TracingTags:        make(map[string]string),
		ReporterSpec:       opt.BuildReporterSpec(),
	}
}
