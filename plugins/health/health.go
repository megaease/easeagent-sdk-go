package Health

import (
	"net/http"

	"github.com/megaease/easemesh/easeagent-sdk-go/plugins"
)

const (
	// Kind is the kind of Health plugin.
	Kind = "Health"
)

// DefaultSpec returns the default spec of Health.
func DefaultSpec() plugins.Spec {
	return Spec{}
}

func init() {
	cons := &plugins.Constructor{
		Kind:         Kind,
		DefaultSpec:  DefaultSpec,
		SystemPlugin: true,
		NewInstance:  New,
	}

	plugins.Register(cons)
}

type (
	// Health is the Health dedicated plugin.
	Health struct {
		spec Spec
	}

	// Spec is the Health spec.
	Spec struct {
		plugins.BaseSpec `json:",inline"`
	}
)

// Validate validates the Health spec.
func (s Spec) Validate() error {
	return nil
}

// New creates a Health plugin.
func New(spec plugins.Spec) (plugins.Plugin, error) {
	h := &Health{
		spec: spec.(Spec),
	}

	return h, nil
}

// HandleAgentRequest handles the agent request.
func (h *Health) HandleAgentRequest(w http.ResponseWriter, r *http.Request) bool {
	switch r.URL.Path {
	case "/health", "/healthz":
		w.WriteHeader(http.StatusOK)
		return true
	default:
		return false
	}
}

// Close closes the Health plugin.
func (h *Health) Close() error {
	return nil
}