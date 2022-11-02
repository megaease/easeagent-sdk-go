package zipkin

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/megaease/easemesh/easeagent-sdk-go/plugins"
	"gopkg.in/yaml.v2"
)

const (
	// Kind is the kind of Zipkin plugin.
	Kind                             = "ZipKin"
	DEFAULT_MEGAEASE_SDK_CONFIG_FILE = "/megaease/sdk/agent.yml"
)

var MEGAEASE_SDK_CONFIG_FILE = os.Getenv("MEGAEASE_SDK_CONFIG_FILE")

// Options is the Zipkin spec.
type Options struct {
	plugins.BaseSpec              `json:",inline"`
	ServiceName                   string  `yaml:"name"`
	TracingEnable                 bool    `default:"true" yaml:"enable"`
	TracingSampleRate             float64 `default:"1" yaml:"tracing.sample.rate" jsonschema:"required,minimum=0,maximum=1"`
	TracingSharedSpans            bool    `default:"true" yaml:"tracing.shared.spans"`
	TracingID128Bit               bool    `default:"true" yaml:"tracing.id128bit"`
	ReporterOutputServer          string  `yaml:"reporter.output.server"`
	ReporterOutputServerTlsEnable bool    `yaml:"reporter.output.server.tls.enable"`
	ReporterOutputServerTlsKey    string  `yaml:"reporter.output.server.tls.key"`
	ReporterOutputServerTlsCert   string  `yaml:"reporter.output.server.tls.cert"`
	ReporterOutputServerTlsCaCert string  `yaml:"reporter.output.server.tls.ca_cert"`
	ReporterTracingSenderUrl      string  `yaml:"reporter.tracing.sender.url"`

	HostPort   string
	HomeDir    string `yaml:"home-dir" long:"home-dir" description:"Path to the home directory."`
	ConfigFile string `yaml:"-" short:"f" long:"config-file" description:"Agent configuration from a file(yaml format), other command line flags will be ignored if specified."`
}

func NewOptions() *Options {
	o := &Options{
		BaseSpec: plugins.BaseSpec{
			NameField: "demo-service",
		},
	}
	var err error
	o.HomeDir, err = filepath.Abs(path.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("failed to identify the full home dir: %v", err)
	}

	return o
}

func LoadGlobalOptions() *Options {
	options := NewOptions()
	err := options.Parse()
	if err != nil {
		log.Panicf("failed to Parse Options path: %v", err)
	}
	options.KindField = Kind
	options.NameField = options.ServiceName
	return options
}

// Validate validates the Zipkin spec.
func (opt Options) Validate() error {
	if opt.ServiceName == "" {
		return fmt.Errorf("name must not be empty")
	}
	return nil
}

func (opt *Options) Parse() error {
	if cfile, err := configFile(); err == nil {
		opt.ConfigFile = cfile
	}
	if opt.ConfigFile == "" {
		return nil
	}
	buff, err := ioutil.ReadFile(opt.ConfigFile)
	if err != nil {
		return fmt.Errorf("read config file failed: %v", err)
	}
	err = yaml.Unmarshal(buff, opt)
	if err != nil {
		return fmt.Errorf("unmarshal config file %s to yaml failed: %v",
			opt.ConfigFile, err)
	}
	return nil
}

func (opt *Options) BuildReporterSpec() *ReporterSpec {
	return &ReporterSpec{
		SpanSpec: &SpanSpec{
			Service: opt.ServiceName,
		},
		SenderUrl: opt.ReporterOutputServer + opt.ReporterTracingSenderUrl,
		TlsEnable: opt.ReporterOutputServerTlsEnable,
		TlsKey:    opt.ReporterOutputServerTlsKey,
		TlsCert:   opt.ReporterOutputServerTlsCert,
		TlsCaCert: opt.ReporterOutputServerTlsCaCert,
	}
}

func (opt *Options) BuildTracingSpec() *TracingSpec {
	return &TracingSpec{
		HostPort:           opt.HostPort,
		ServiceName:        opt.ServiceName,
		TracingEnable:      opt.TracingEnable,
		TracingSampleRate:  opt.TracingSampleRate,
		TracingSharedSpans: opt.TracingSharedSpans,
		TracingID128Bit:    opt.TracingID128Bit,
		TracingTags:        make(map[string]string),
		ReporterSpec:       opt.BuildReporterSpec(),
	}
}

func configFile() (string, error) {
	cfile := MEGAEASE_SDK_CONFIG_FILE
	if cfile == "" {
		cfile = path.Join(path.Dir(os.Args[0]), "agent.yml")
		_, err := os.Stat(cfile)
		if err != nil {
			log.Printf("file path: %v", err)
			cfile = ""
		}
	}
	if cfile == "" {
		_, err := os.Stat(DEFAULT_MEGAEASE_SDK_CONFIG_FILE)
		if err == nil {
			cfile = DEFAULT_MEGAEASE_SDK_CONFIG_FILE
		}
	}
	log.Printf("cfile: %s\n", cfile)
	_, err := os.Stat(cfile)
	if err != nil {
		log.Printf("failed to get config file path: %v", err)
		return "", err
	}
	return cfile, nil
}
