package vault

import (
	"fmt"
	"sync"
	"time"
)

func maturityKey(mount, path string) string {
	return mount + "/" + path
}

// SecretMaturityRegistry stores maturity assignments keyed by mount+path.
type SecretMaturityRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretMaturity
}

// NewSecretMaturityRegistry returns an initialised registry.
func NewSecretMaturityRegistry() *SecretMaturityRegistry {
	return &SecretMaturityRegistry{
		entries: make(map[string]SecretMaturity),
	}
}

// Set validates and stores a maturity record, setting AssignedAt if zero.
func (r *SecretMaturityRegistry) Set(m SecretMaturity) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if m.AssignedAt.IsZero() {
		m.AssignedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[maturityKey(m.Mount, m.Path)] = m
	return nil
}

// Get retrieves the maturity record for the given mount and path.
func (r *SecretMaturityRegistry) Get(mount, path string) (SecretMaturity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.entries[maturityKey(mount, path)]
	if !ok {
		return SecretMaturity{}, fmt.Errorf("maturity: no record for %s/%s", mount, path)
	}
	return m, nil
}

// Remove deletes the maturity record for the given mount and path.
func (r *SecretMaturityRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, maturityKey(mount, path))
}

// All returns a snapshot of every stored maturity record.
func (r *SecretMaturityRegistry) All() []SecretMaturity {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretMaturity, 0, len(r.entries))
	for _, v := range r.entries {
		out = append(out, v)
	}
	return out
}
