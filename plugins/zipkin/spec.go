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

package zipkin

import (
	"fmt"

	"github.com/megaease/easeagent-sdk-go/plugins"
)

const (
	// Kind is the kind of Zipkin plugin.
	Kind = "Zipkin"
	// Name is the name of Zipkin plugin.
	Name = "Zipkin"
)

type (
	// Spec is the Zipkin spec.
	Spec struct {
		plugins.BaseSpec `json:",inline"`

		OutputServerURL string `json:"reporter.output.server"`

		EnableTLS bool   `json:"reporter.output.server.tls.enable"`
		TLSKey    string `json:"reporter.output.server.tls.key"`
		TLSCert   string `json:"reporter.output.server.tls.cert"`
		TLSCaCert string `json:"reporter.output.server.tls.caCert"`

		EnableBasicAuth bool   `json:"reporter.output.server.auth.enable"`
		Username        string `json:"reporter.output.server.auth.username"`
		Password        string `json:"reporter.output.server.auth.password"`

		ServiceName   string            `json:"serviceName"`
		TracingType   string            `json:"tracing.type"`
		LocalHostport string            `json:"-"`
		Tags          map[string]string `json:"-"`

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
			NameField: Name,
		},
		OutputServerURL: "https://127.0.0.1:8080/report",

		EnableTLS: false,

		EnableBasicAuth: false,

		EnableTracing: true,
		ServiceName:   "default-service",
		TracingType:   "log-tracing",
		LocalHostport: "127.0.0.1:80",

		SampleRate:  1,
		SharedSpans: true,
		ID128Bit:    false,
		Username:    "",
		Password:    "",
	}
}

// NewConsoleReportSpec new a Console Reporter Spec
func NewConsoleReportSpec(localHostPort string) Spec {
	spec := DefaultSpec().(Spec)
	spec.OutputServerURL = "" // report to log when output server is ""
	spec.LocalHostport = localHostPort
	return spec
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
