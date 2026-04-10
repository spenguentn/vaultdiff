package vault

import "fmt"

// EngineRegistry holds a collection of known secrets engines indexed by mount.
type EngineRegistry struct {
	engines map[string]SecretsEngine
}

// NewEngineRegistry returns an initialised, empty EngineRegistry.
func NewEngineRegistry() *EngineRegistry {
	return &EngineRegistry{
		engines: make(map[string]SecretsEngine),
	}
}

// Register adds or replaces a SecretsEngine in the registry.
// Returns an error if the engine fails validation.
func (r *EngineRegistry) Register(e SecretsEngine) error {
	if err := e.Validate(); err != nil {
		return fmt.Errorf("register engine: %w", err)
	}
	r.engines[e.Mount] = e
	return nil
}

// Get retrieves a SecretsEngine by its mount path.
// The second return value is false when the mount is not registered.
func (r *EngineRegistry) Get(mount string) (SecretsEngine, bool) {
	e, ok := r.engines[mount]
	return e, ok
}

// Remove deletes the engine with the given mount path from the registry.
func (r *EngineRegistry) Remove(mount string) {
	delete(r.engines, mount)
}

// All returns a slice of every registered SecretsEngine.
func (r *EngineRegistry) All() []SecretsEngine {
	out := make([]SecretsEngine, 0, len(r.engines))
	for _, e := range r.engines {
		out = append(out, e)
	}
	return out
}

// Count returns the number of registered engines.
func (r *EngineRegistry) Count() int {
	return len(r.engines)
}
