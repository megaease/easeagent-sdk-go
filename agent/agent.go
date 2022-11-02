package agent

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/megaease/easemesh/easeagent-sdk-go/plugins"
	"github.com/megaease/easemesh/easeagent-sdk-go/plugins/zipkin"
	"golang.org/x/exp/maps"
)

const (
	// DefualtAgentPort is the default listeaning port for agent.
	// https://github.com/megaease/easemesh/blob/main/docs/sidecar-protocol.md#easemesh-traffic-hosting
	DefualtAgentPort = 9900

	// DefaultAgentType is the default global agent type.
	DefaultAgentType = "GoSDK"

	// AgentVersion is the current version of Agent.
	AgentVersion = "v0.1.0"
)

var (
	agentAddr = fmt.Sprintf(":%d", DefualtAgentPort)

	// DefaultAgentConfig is the global default agent.
	DefaultAgentConfig = &Config{
		Address: agentAddr,
		Plugins: []plugins.Spec{
			zipkin.DefaultSpec(),
		},
	}

	// DefaultAgent is the default global agent.
	DefaultAgent *Agent
)

func init() {
	agent, err := New(DefaultAgentConfig)
	if err != nil {
		panic(err)
	}
	DefaultAgent = agent
}

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

	return agent, nil
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

func (a *Agent) WrapHttpRequest(parent context.Context, req *http.Request) *http.Request {
	for _, plug := range a.plugins {
		wrapper, ok := plug.(plugins.UserClientRequestWrapper)
		if !ok {
			continue
		}
		req = wrapper.WrapUserClientRequest(parent, req)
	}

	return req
}

// ServeDefaultAgent just runs global default agent in HTTP server,
// please notice it prints logs if the server failed listening.
// The caller must call it to activate default agent.
func ServeDefaultAgent() {
	go func() {
		err := http.ListenAndServe(agentAddr, DefaultAgent)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("easemesh agent listen %s failed: %v", agentAddr, err)
		}
	}()
}

// ServeAgent just runs the given agent in HTTP server,
// please notice it prints logs if the server failed listening.
func ServeAgent(agent *Agent) {
	go func() {
		err := http.ListenAndServe(agentAddr, agent)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("easemesh agent listen %s failed: %v", agentAddr, err)
		}
	}()
}
