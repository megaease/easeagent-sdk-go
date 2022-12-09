package agent

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/megaease/easeagent-sdk-go/plugins"
	"github.com/megaease/easeagent-sdk-go/plugins/zipkin"
	"gopkg.in/yaml.v2"
)

// ConfigOption allows for functional options to adjust behavior and payload of
// the Config to be created with agent.Agent().
type ConfigOption func(c *Config)

// NewWithOptions returns a new Agent.
func NewWithOptions(options ...ConfigOption) (*Agent, error) {
	config := &Config{
		Plugins: make([]plugins.Spec, 0),
	}
	for _, option := range options {
		option(config)
	}
	return New(config)
}

// WithAddress Sets address to Config.
// @param  address string agent http api address
// @return ConfigOption
func WithAddress(address string) ConfigOption {
	return func(c *Config) {
		c.Address = address
	}
}

// WithSpec Append spec the Agent Plugin Spec.
// @param  spec plugins.Spec
// @return ConfigOption
func WithSpec(spec plugins.Spec) ConfigOption {
	return func(c *Config) {
		c.Plugins = append(c.Plugins, spec)
	}
}

// WithZipkinYaml Append zipkin spec load from yaml file to the Agent Plugin Spec.
//  			  sets host and port of the tracer Span.localEndpoint.
// @param  yamlFile string yaml file path. use yamlFile="" is Console Reporter for tracing.
// @param  localHostPort string host and port of the tracer Span.localEndpoint.
// 								By default, use localHostPort="" is not set host and port of Span.localEndpoint.
// @return ConfigOption
func WithZipkinYaml(yamlFile string, localHostPort string) ConfigOption {
	return func(c *Config) {
		var spec zipkin.Spec
		var body map[string]interface{}
		if yamlFile == "" {
			log.Printf("yamlFile was '', use default Console Reporter for tracing.")
			spec = zipkin.NewConsoleReportSpec(localHostPort)
		} else if buff, err := ioutil.ReadFile(yamlFile); err != nil {
			log.Printf("read config file:%s failed: %v, use default Console Reporter for tracing.", yamlFile, err)
			spec = zipkin.NewConsoleReportSpec(localHostPort)
		} else if err = yaml.Unmarshal(buff, &body); err != nil {
			log.Printf("unmarshal yaml file %s to map failed: %v, use default Console Reporter for tracing.", yamlFile, err)
			spec = zipkin.NewConsoleReportSpec(localHostPort)
		} else if bodyJSON, err := json.Marshal(body); err != nil {
			log.Printf("marshal yaml file %s to json failed: %v, use default Console Reporter for tracing.", yamlFile, err)
			spec = zipkin.NewConsoleReportSpec(localHostPort)
		} else if err = json.Unmarshal(bodyJSON, &spec); err != nil {
			log.Printf("unmarshal %s to %T failed: %v, use default Console Reporter for tracing.", bodyJSON, spec, err)
			spec = zipkin.NewConsoleReportSpec(localHostPort)
		} else {
			spec.KindField = zipkin.Kind
			spec.NameField = zipkin.NAME
			spec.LocalHostport = localHostPort
		}
		c.Plugins = append(c.Plugins, spec)
	}
}
