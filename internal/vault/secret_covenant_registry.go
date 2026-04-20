package vault

import (
	"fmt"
	"sync"
	"time"
)

func covenantKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretCovenantRegistry stores and retrieves covenants keyed by mount+path.
type SecretCovenantRegistry struct {
	mu    sync.RWMutex
	store map[string]*SecretCovenant
}

// NewSecretCovenantRegistry returns an initialised SecretCovenantRegistry.
func NewSecretCovenantRegistry() *SecretCovenantRegistry {
	return &SecretCovenantRegistry{
		store: make(map[string]*SecretCovenant),
	}
}

// Set validates and stores a covenant, setting CreatedAt when it is zero.
func (r *SecretCovenantRegistry) Set(c *SecretCovenant) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[covenantKey(c.Mount, c.Path)] = c
	return nil
}

// Get returns the covenant for the given mount and path, or an error if not found.
func (r *SecretCovenantRegistry) Get(mount, path string) (*SecretCovenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.store[covenantKey(mount, path)]
	if !ok {
		return nil, fmt.Errorf("covenant: no entry for %s/%s", mount, path)
	}
	return c, nil
}

// Remove deletes the covenant for the given mount and path.
func (r *SecretCovenantRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, covenantKey(mount, path))
}

// All returns a snapshot of every covenant currently registered.
func (r *SecretCovenantRegistry) All() []*SecretCovenant {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*SecretCovenant, 0, len(r.store))
	for _, c := range r.store {
		out = append(out, c)
	}
	return out
}
