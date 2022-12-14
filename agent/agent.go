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
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/megaease/easeagent-sdk-go/plugins"
	"golang.org/x/exp/maps"
)

type (
	// Agent is the agent entry.
	Agent struct {
		config  *Config
		plugins []plugins.Plugin
	}

	// Config is the Agent config.
	Config struct {
		Address string         `json:"address"`
		Plugins []plugins.Spec `json:"plugins"`
	}

	// HandlerWrapper is the HTTP handler wrapper.
	HandlerWrapper struct {
		handlerFunc http.HandlerFunc
	}
)

// New creates an agent.
func New(config *Config) (*Agent, error) {
	agent := &Agent{
		config: config,
	}

	systemConstructors := plugins.SystemConstructors()
	var plugs []plugins.Plugin
	for i, spec := range config.Plugins {
		plug, err := plugins.New(spec)
		if err != nil {
			return nil, fmt.Errorf("failed to create No.%d plugin: %v", i+1, err)
		}
		plugs = append(plugs, plug)
		delete(systemConstructors, spec.Kind())
	}

	// NOTE: For consistensy.

	systemCons := maps.Values(systemConstructors)
	sort.Sort(plugins.ConstructorsByKind(systemCons))

	for _, cons := range systemCons {
		plug, err := cons.NewInstance(cons.DefaultSpec())
		if err != nil {
			return nil, fmt.Errorf("failed to create system plugin %s: %v", cons.Kind, err)
		}
		plugs = append(plugs, plug)
	}
	agent.plugins = plugs

	go func() {
		err := http.ListenAndServe(config.Address, agent)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("easemesh agent listen %s failed: %v", config.Address, err)
		}
	}()

	return agent, nil
}

// GetPlugin gets Plugin by name
func (a *Agent) GetPlugin(name string) plugins.Plugin {
	for _, plug := range a.plugins {
		if plug.Name() == name {
			return plug
		}
	}
	return nil
}

// ServeHTTP invokes every plugin which is http.Handler to handle the request.
// NOTE: If the request is not matched for your plugin, don't do anything.
func (a *Agent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, plug := range a.plugins {
		handler, ok := plug.(plugins.AgentHandler)
		if !ok {
			continue
		}

		handled := handler.HandleAgentRequest(w, r)
		if handled {
			return
		}
	}
}

// WrapUserHandlerFunc wraps the handlerFunc with the wrap functions from every plugin.
func (a *Agent) WrapUserHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	for _, plug := range a.plugins {
		wrapper, ok := plug.(plugins.UserHandlerFuncWrapper)
		if !ok {
			continue
		}

		handlerFunc = wrapper.WrapUserHandlerFunc(handlerFunc)
	}

	return handlerFunc
}

// WrapUserHandler wraps the handler with the wrap functions from every enabled plugin.
func (a *Agent) WrapUserHandler(handler http.Handler) http.Handler {
	return &HandlerWrapper{
		handlerFunc: a.WrapUserHandlerFunc(handler.ServeHTTP),
	}
}

func (h *HandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handlerFunc(w, r)
}

// WrapUserClient wraps the client with the wrap functions from every enabled plugin.
func (a *Agent) WrapUserClient(httpDoer plugins.HTTPDoer) plugins.HTTPDoer {
	for _, plug := range a.plugins {
		wrapper, ok := plug.(plugins.UserClientWrapper)
		if !ok {
			continue
		}

		httpDoer = wrapper.WrapUserClient(httpDoer)
	}

	return httpDoer
}

// WrapHTTPRequest wraps the request with the parent conext from every enabled plugin.
func (a *Agent) WrapHTTPRequest(parent context.Context, req *http.Request) *http.Request {
	for _, plug := range a.plugins {
		wrapper, ok := plug.(plugins.UserClientRequestWrapper)
		if !ok {
			continue
		}
		req = wrapper.WrapUserClientRequest(parent, req)
	}

	return req
}
