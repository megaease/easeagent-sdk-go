package agent

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	DEFAULT_MEGAEASE_SDK_CONFIG_FILE = "/megaease/sdk/agent.yml"
	MEGAEASE_SDK_CONFIG_FILE         = os.Getenv("MEGAEASE_SDK_CONFIG_FILE")
	GlobalOptions                    *Options
)

type Options struct {
	Name                          string  `yaml:"name"`
	System                        string  `yaml:"system"`
	SampleRate                    float64 `yaml:"sampleRate" jsonschema:"required,minimum=0,maximum=1"`
	ReporterOutputServer          string  `yaml:"reporter.output.server"`
	ReporterOutputServerTlsEnable bool    `yaml:"reporter.output.server.tls.enable"`
	ReporterOutputServerTlsKey    string  `yaml:"reporter.output.server.tls.key"`
	ReporterOutputServerTlsCert   string  `yaml:"reporter.output.server.tls.cert"`
	ReporterOutputServerTlsCaCert string  `yaml:"reporter.output.server.tls.ca_cert"`
	ReporterTracingSenderUrl      string  `yaml:"reporter.tracing.sender.url"`

	HomeDir    string `yaml:"home-dir" long:"home-dir" description:"Path to the home directory."`
	ConfigFile string `yaml:"-" short:"f" long:"config-file" description:"Agent configuration from a file(yaml format), other command line flags will be ignored if specified."`
}

func New() *Options {
	o := &Options{
		Name:   "demo-service",
		System: "demo-system",
	}
	var err error
	o.HomeDir, err = filepath.Abs(path.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("failed to identify the full home dir: %v", err)
	}

	return o
}

func InitGlobalOptions() {
	options := New()
	err := options.Parse()
	if err != nil {
		log.Panicf("failed to Parse Options path: %v", err)
	}
	GlobalOptions = options
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
