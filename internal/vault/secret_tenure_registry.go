package vault

import (
	"fmt"
	"sync"
	"time"
)

func tenureKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretTenureRegistry stores tenure records keyed by mount+path.
type SecretTenureRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretTenure
}

// NewSecretTenureRegistry returns an initialised SecretTenureRegistry.
func NewSecretTenureRegistry() *SecretTenureRegistry {
	return &SecretTenureRegistry{
		entries: make(map[string]SecretTenure),
	}
}

// Set computes and stores the tenure for the given secret.
func (r *SecretTenureRegistry) Set(mount, path string, createdAt time.Time) error {
	t := ComputeTenure(mount, path, createdAt)
	if err := t.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[tenureKey(mount, path)] = t
	return nil
}

// Get returns the stored tenure for a secret, or false if not found.
func (r *SecretTenureRegistry) Get(mount, path string) (SecretTenure, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.entries[tenureKey(mount, path)]
	return t, ok
}

// Remove deletes the tenure record for the given secret.
func (r *SecretTenureRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, tenureKey(mount, path))
}

// All returns a snapshot of every tenure record currently stored.
func (r *SecretTenureRegistry) All() []SecretTenure {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretTenure, 0, len(r.entries))
	for _, v := range r.entries {
		out = append(out, v)
	}
	return out
}
