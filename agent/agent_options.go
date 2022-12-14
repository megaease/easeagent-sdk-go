/**
 * Copyright 2022 MegaEase
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *     http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/megaease/easeagent-sdk-go/plugins"
	"github.com/megaease/easeagent-sdk-go/plugins/easemesh"
	"github.com/megaease/easeagent-sdk-go/plugins/health"
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

// WithZipkinYAML Append easemesh spec load from yaml file to the Agent Plugin Spec.
// @param  yamlFile string yaml file path. use yamlFile="" is use easemesh.DefaultSpec() for
// @return ConfigOption
func WithEaseMeshYAML(yamlFile string) ConfigOption {
	return func(c *Config) {
		var easeMeshSpec easemesh.Spec
		var spec plugins.Spec
		if bodyJSON, err := yamlToJSON(yamlFile); err != nil {
			log.Printf("yaml to json failed: %v, use default easemesh spec", err)
			spec = easemesh.DefaultSpec()
		} else if err = json.Unmarshal(bodyJSON, &easeMeshSpec); err != nil {
			log.Printf("unmarshal %s to %T failed: %v, use default easemesh spec", bodyJSON, spec, err)
			spec = easemesh.DefaultSpec()
		} else {
			easeMeshSpec.KindField = easemesh.Kind
			easeMeshSpec.NameField = easemesh.Name
			spec = easeMeshSpec
		}
		c.Plugins = append(c.Plugins, spec)
	}
}

// WithZipkinYAML Append zipkin spec load from yaml file to the Agent Plugin Spec.
//  			  sets host and port of the tracer Span.localEndpoint.
// @param  yamlFile string yaml file path. use yamlFile="" is Console Reporter for tracing.
// @param  localHostPort string host and port of the tracer Span.localEndpoint.
// 								By default, use localHostPort="" is not sets host and port of Span.localEndpoint.
// @return ConfigOption
func WithZipkinYAML(yamlFile string, localHostPort string) ConfigOption {
	return func(c *Config) {
		var spec zipkin.Spec
		if bodyJSON, err := yamlToJSON(yamlFile); err != nil {
			log.Printf("yaml to json failed: %v, use default Console Reporter for tracing", err)
			spec = zipkin.NewConsoleReportSpec(localHostPort)
		} else if err = json.Unmarshal(bodyJSON, &spec); err != nil {
			log.Printf("unmarshal %s to %T failed: %v, use default Console Reporter for tracing.", bodyJSON, spec, err)
			spec = zipkin.NewConsoleReportSpec(localHostPort)
		} else {
			spec.KindField = zipkin.Kind
			spec.NameField = zipkin.Name
			spec.LocalHostport = localHostPort
		}
		c.Plugins = append(c.Plugins, spec)
	}
}

// WithYAML sets address, Append health, easemesh and zipkin spec load from yaml file to the Agent Plugin Spec.
// @param  yamlFile string yaml file path. use yamlFile="" is use easemesh.DefaultSpec() and Console Reporter for tracing.
// @param  localHostPort string host and port of the tracer Span.localEndpoint.
// 								By default, use localHostPort="" is not sets host and port of Span.localEndpoint.
// @return ConfigOption
func WithYAML(yamlFile string, localHostPort string) ConfigOption {
	return func(c *Config) {
		bodyJSON, err := yamlToJSON(yamlFile)
		if err == nil {
			err = json.Unmarshal(bodyJSON, &c)
			if err != nil {
				log.Printf("unmarshal %s to %T failed: %v, can't load base config", bodyJSON, c, err)
			}
		}
		c.Plugins = append(c.Plugins, health.DefaultSpec())
		WithEaseMeshYAML(yamlFile)(c)
		WithZipkinYAML(yamlFile, localHostPort)(c)
	}
}

func yamlToJSON(yamlFile string) ([]byte, error) {
	var body map[string]interface{}
	if yamlFile == "" {
		return nil, fmt.Errorf("yamlFile was ''")
	} else if buff, err := ioutil.ReadFile(yamlFile); err != nil {
		return nil, fmt.Errorf("read config file:%s failed: %v", yamlFile, err)
	} else if err = yaml.Unmarshal(buff, &body); err != nil {
		return nil, fmt.Errorf("unmarshal yaml file %s to map failed: %v", yamlFile, err)
	} else if bodyJSON, err := json.Marshal(body); err != nil {
		return nil, fmt.Errorf("marshal yaml file %s to json failed: %v", yamlFile, err)
	} else {
		return bodyJSON, err
	}
}

// WithZipkinTags sets tags of Zipkin Plugin Spec
// @param tags the Span tags of tracing
func WithZipkinTags(tags map[string]string) ConfigOption {
	return func(c *Config) {
		for i := 0; i < len(c.Plugins); i++ {
			plugin := c.Plugins[i]
			if zipkinSpec, ok := plugin.(zipkin.Spec); ok {
				zipkinSpec.Tags = tags
			}
		}
	}
}
