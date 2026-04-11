package vault

import (
	"fmt"
	"sync"
)

// SecretWatchRegistry manages a collection of active secret watch configs.
type SecretWatchRegistry struct {
	mu      sync.RWMutex
	watches map[string]SecretWatchConfig
}

// NewSecretWatchRegistry returns an initialised SecretWatchRegistry.
func NewSecretWatchRegistry() *SecretWatchRegistry {
	return &SecretWatchRegistry{
		watches: make(map[string]SecretWatchConfig),
	}
}

func watchKey(mount, path string) string {
	return fmt.Sprintf("%s::%s", mount, path)
}

// Register adds a watch config to the registry after validation.
func (r *SecretWatchRegistry) Register(cfg SecretWatchConfig) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.watches[watchKey(cfg.Mount, cfg.Path)] = cfg
	return nil
}

// Get retrieves a watch config by mount and path.
func (r *SecretWatchRegistry) Get(mount, path string) (SecretWatchConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cfg, ok := r.watches[watchKey(mount, path)]
	return cfg, ok
}

// Remove deletes a watch config from the registry.
func (r *SecretWatchRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.watches, watchKey(mount, path))
}

// List returns all registered watch configs.
func (r *SecretWatchRegistry) List() []SecretWatchConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretWatchConfig, 0, len(r.watches))
	for _, cfg := range r.watches {
		out = append(out, cfg)
	}
	return out
}
