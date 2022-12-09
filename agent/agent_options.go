package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/megaease/easeagent-sdk-go/plugins"
	"github.com/megaease/easeagent-sdk-go/plugins/zipkin"
	"gopkg.in/yaml.v2"
)

// ConfigOption allows for functional options to adjust behavior and payload of
// the Config to be created with agent.Agent().
type ConfigOption func(c *Config) error

// NewAgent returns a new Agent.
func NewWithOptions(options ...ConfigOption) (*Agent, error) {
	config := &Config{
		Plugins: make([]plugins.Spec, 0),
	}
	for _, option := range options {
		if err := option(config); err != nil {
			return nil, err
		}
	}
	return New(config)
}

func WithAddress(address string) ConfigOption {
	return func(c *Config) error {
		c.Address = address
		return nil
	}
}

func WithSpec(spec plugins.Spec) ConfigOption {
	return func(c *Config) error {
		c.Plugins = append(c.Plugins, spec)
		return nil
	}
}

// WithZipkinYaml append zipkin spec load from yaml file to the Agent Plugin Spec, sets the local endpoint host and port of the tracer.
func WithZipkinYaml(yamlFile string, localHostPort string) ConfigOption {
	return func(c *Config) error {
		var spec zipkin.Spec
		if yamlFile == "" {
			spec = zipkin.DefaultSpec().(zipkin.Spec)
			spec.OutputServerURL = "" // report to log when output server is ""
			spec.LocalHostport = localHostPort
		} else {
			buff, err := ioutil.ReadFile(yamlFile)
			if err != nil {
				return fmt.Errorf("read config file :%s failed: %v", yamlFile, err)
			}
			var body map[string]interface{}
			err = yaml.Unmarshal(buff, &body)
			if err != nil {
				return fmt.Errorf("unmarshal yaml file %s to map failed: %v",
					yamlFile, err)
			}

			bodyJson, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("marshal yaml file %s to json failed: %v",
					yamlFile, err)
			}
			err = json.Unmarshal(bodyJson, &spec)
			if err != nil {
				return fmt.Errorf("unmarshal %s to %T failed: %v", bodyJson, spec, err)
			}
			spec.KindField = zipkin.Kind
			spec.NameField = zipkin.NAME
			spec.LocalHostport = localHostPort
		}
		c.Plugins = append(c.Plugins, spec)

		return nil
	}
}
