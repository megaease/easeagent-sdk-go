package easemesh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/megaease/easeagent-sdk-go/plugins"
	"golang.org/x/exp/slices"
)

const (
	// Kind is the kind of EaseMesh plugin.
	Kind = "EaseMesh"

	agentType    = "GoSDK"
	agentVersion = "v0.1.0"
)

// DefaultSpec returns the default spec of EaseMesh.
func DefaultSpec() plugins.Spec {
	return Spec{}
}

func init() {
	cons := &plugins.Constructor{
		Kind:         Kind,
		DefaultSpec:  DefaultSpec,
		SystemPlugin: false,
		NewInstance:  New,
	}

	plugins.Register(cons)
}

type (
	// EaseMesh is the EaseMesh dedicated plugin.
	EaseMesh struct {
		spec Spec

		agentInfo []byte
		headers   atomic.Value // type: []string
	}

	// Spec is the EaseMesh spec.
	Spec struct {
		plugins.BaseSpec `json:",inline"`
	}

	// AgentInfo stores agent information.
	AgentInfo struct {
		Type    string `json:"type"`
		Version string `json:"version"`
	}

	// AgentConfig is the config pushed to agent.
	AgentConfig struct {
		Headers string `json:"easeagent.progress.forwarded.headers"`
	}
)

// Validate validates the EaseMesh spec.
func (s Spec) Validate() error {
	return nil
}

// New creates a EaseMesh plugin.
func New(spec plugins.Spec) (plugins.Plugin, error) {
	agentInfo := &AgentInfo{
		Type:    agentType,
		Version: agentVersion,
	}

	buff, err := json.Marshal(agentInfo)
	if err != nil {
		return nil, fmt.Errorf("marshal agent info failed: %v", err)
	}

	mesh := &EaseMesh{
		agentInfo: buff,
		spec:      spec.(Spec),
	}

	mesh.headers.Store([]string{})

	return mesh, nil
}

// HandleAgentRequest handles the agent request.
func (mesh *EaseMesh) HandleAgentRequest(w http.ResponseWriter, r *http.Request) bool {
	switch r.URL.Path {
	case "/config":
		mesh.handleConfig(w, r)
		return true
	case "/agent-info":
		mesh.handleAgentInfo(w, r)
		return true
	default:
		return false
	}
}

// WrapUserHandlerFunc wraps the user handler function.
func (mesh *EaseMesh) WrapUserHandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys := mesh.headers.Load().([]string)
		for k, v := range r.Header.Clone() {
			if slices.Contains(keys, k) {
				w.Header()[k] = v
			}
		}

		fn(w, r)

		// NOTE: Copying headers after fn it might not take effect,
		// in the case of fn invoking w.WriteHeader.
	}
}

func (mesh *EaseMesh) handleConfig(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read config body failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	config := &AgentConfig{}
	err = json.Unmarshal(body, config)
	if err != nil {
		log.Printf("unmarshal config body failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	headers := strings.Split(config.Headers, ",")
	mesh.headers.Store(headers)
}

func (mesh *EaseMesh) handleAgentInfo(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(mesh.agentInfo)
	if err != nil {
		log.Printf("marshal agent info failed: %v", err)
		w.WriteHeader(500)
	}

	w.Write(data)
}

// Close closes the EaseMesh plugin.
func (mesh *EaseMesh) Close() error {
	return nil
}
