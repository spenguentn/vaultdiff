package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// EnvironmentBinding associates a secret path with a specific environment,
// allowing secrets to be scoped and validated per deployment target.
type EnvironmentBinding struct {
	Mount       string    `json:"mount"`
	Path        string    `json:"path"`
	Environment string    `json:"environment"`
	BoundBy     string    `json:"bound_by"`
	BoundAt     time.Time `json:"bound_at"`
	Required    bool      `json:"required"`
}

// FullPath returns the canonical mount+path identifier.
func (b EnvironmentBinding) FullPath() string {
	return fmt.Sprintf("%s/%s", strings.Trim(b.Mount, "/"), strings.Trim(b.Path, "/"))
}

// Validate checks that the binding contains required fields.
func (b EnvironmentBinding) Validate() error {
	if b.Mount == "" {
		return errors.New("environment binding: mount is required")
	}
	if b.Path == "" {
		return errors.New("environment binding: path is required")
	}
	if b.Environment == "" {
		return errors.New("environment binding: environment is required")
	}
	if b.BoundBy == "" {
		return errors.New("environment binding: bound_by is required")
	}
	return nil
}

// bindingKey returns a unique registry key for a binding.
func bindingKey(mount, path, environment string) string {
	return fmt.Sprintf("%s|%s|%s",
		strings.Trim(mount, "/"),
		strings.Trim(path, "/"),
		strings.ToLower(environment),
	)
}

// NewEnvironmentBindingRegistry creates an empty registry.
func NewEnvironmentBindingRegistry() *EnvironmentBindingRegistry {
	return &EnvironmentBindingRegistry{
		bindings: make(map[string]EnvironmentBinding),
	}
}

// EnvironmentBindingRegistry stores and retrieves environment bindings.
type EnvironmentBindingRegistry struct {
	bindings map[string]EnvironmentBinding
}

// Bind registers a binding, setting BoundAt if not already set.
func (r *EnvironmentBindingRegistry) Bind(b EnvironmentBinding) error {
	if err := b.Validate(); err != nil {
		return err
	}
	if b.BoundAt.IsZero() {
		b.BoundAt = time.Now().UTC()
	}
	r.bindings[bindingKey(b.Mount, b.Path, b.Environment)] = b
	return nil
}

// Get retrieves a binding by mount, path, and environment.
func (r *EnvironmentBindingRegistry) Get(mount, path, environment string) (EnvironmentBinding, bool) {
	b, ok := r.bindings[bindingKey(mount, path, environment)]
	return b, ok
}

// Remove deletes a binding from the registry.
func (r *EnvironmentBindingRegistry) Remove(mount, path, environment string) {
	delete(r.bindings, bindingKey(mount, path, environment))
}

// All returns all registered bindings.
func (r *EnvironmentBindingRegistry) All() []EnvironmentBinding {
	out := make([]EnvironmentBinding, 0, len(r.bindings))
	for _, b := range r.bindings {
		out = append(out, b)
	}
	return out
}
