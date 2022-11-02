package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	// Plugin is the plugin interface.
	Plugin interface {
		Close() error
	}

	// AgentHandler is the handler to handle requests from agent port (default 9900).
	AgentHandler interface {
		// If the plugin doesn't handle the request, it should return false.
		HandleAgentRequest(w http.ResponseWriter, r *http.Request) bool
	}

	// UserHandlerFuncWrapper wraps the user HandleFunc.
	UserHandlerFuncWrapper interface {
		// If the plugin doesn't wrap the user handler, it should not implement this method.
		WrapUserHandlerFunc(fn http.HandlerFunc) http.HandlerFunc
	}

	// UserClientWrapper wraps the user client.
	UserClientWrapper interface {
		WrapUserClient(client HTTPDoer) HTTPDoer
	}

	// UserClientRequestWrapper wraps the user client.
	UserClientRequestWrapper interface {
		WrapUserClientRequest(parent context.Context, req *http.Request) *http.Request
	}

	// HTTPDoer is the interface to do HTTP request.
	HTTPDoer interface {
		Do(req *http.Request) (*http.Response, error)
	}

	// Constructor contains metadata and contruction resources of a kind of plugin.
	Constructor struct {
		// Kind is the kind of the plugin.
		Kind string

		// DefaultSpec returns a default spec for the plugin, with default values.
		// The function should always return a new spec copy, since it will be
		// modified and used by next moves.
		//
		// NOTE1: The DefaultSpec must returns complete spec (including name) for system plugins.
		DefaultSpec func() Spec

		// SystemPlugin being true value indicates the plugin is system-level,
		// which means it will always be loaded even if not specified in the config.
		SystemPlugin bool

		// NewInstance creates a new plugin instance for the kind.
		// NOTE: The spec has the same type with the one which DefaultSpec returns.
		NewInstance func(spec Spec) (Plugin, error)
	}

	// Spec is the common interface of filter specs
	Spec interface {
		// Name returns name.
		Name() string

		// Kind returns kind.
		Kind() string

		// Validate validates the spec.
		Validate() error
	}

	// BaseSpec is the base spec for all plugins.
	BaseSpec struct {
		NameField string `json:"name"`
		KindField string `json:"kind"`
	}
)

// Name returns name.
func (s BaseSpec) Name() string { return s.NameField }

// Kind returns kind.
func (s BaseSpec) Kind() string { return s.KindField }

// constructors is the registry for plugins.
var constructors = map[string]*Constructor{}

// Register registers a filter kind.
func Register(cons *Constructor) {
	if cons.Kind == "" {
		panic(fmt.Errorf("%T: empty kind", cons))
	}

	if consExisted := constructors[cons.Kind]; consExisted != nil {
		msgFmt := "%T and %T got same name: %s"
		panic(fmt.Errorf(msgFmt, cons, consExisted, cons.Kind))
	}

	constructors[cons.Kind] = cons
}

// NewFromJSON creates a plugin instance according to the JSON spec.
func NewFromJSON(specJSON []byte) (Plugin, error) {
	var baseSpec BaseSpec
	if err := json.Unmarshal(specJSON, &baseSpec); err != nil {
		return nil, fmt.Errorf("unmarshal %s to %T failed: %v", specJSON, baseSpec, err)
	}

	cons, existed := constructors[baseSpec.KindField]
	if !existed {
		return nil, fmt.Errorf("plugin kind %s not found", baseSpec.KindField)
	}

	spec := cons.DefaultSpec()

	if err := json.Unmarshal(specJSON, &spec); err != nil {
		return nil, fmt.Errorf("unmarshal %s to %T failed: %v", specJSON, spec, err)
	}

	return New(spec)
}

// New creates a plugin instance according to the spec.
func New(spec Spec) (Plugin, error) {
	err := spec.Validate()
	if err != nil {
		return nil, fmt.Errorf("validate %T failed: %v", spec, err)
	}

	cons := constructors[spec.Kind()]
	if cons == nil {
		return nil, fmt.Errorf("plugin kind %s not found", spec.Kind())
	}

	if spec.Name() == "" {
		return nil, fmt.Errorf("plugin %s got empty name", spec.Kind())
	}

	instance, err := cons.NewInstance(spec)
	if err != nil {
		return nil, fmt.Errorf("new plugin %s/%s failed: %v", spec.Kind(), spec.Name(), err)
	}

	return instance, nil
}

// ConstructorsByKind is a slice of Constructor, which is sortable by Kind.
type ConstructorsByKind []*Constructor

func (c ConstructorsByKind) Len() int           { return len(c) }
func (c ConstructorsByKind) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ConstructorsByKind) Less(i, j int) bool { return c[i].Kind < c[j].Kind }

// SystemConstructors returns all system plugins.
func SystemConstructors() map[string]*Constructor {
	result := map[string]*Constructor{}
	for _, cons := range constructors {
		if cons.SystemPlugin {
			result[cons.Kind] = cons
		}
	}

	return result
}
