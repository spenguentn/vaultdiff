package vault

import (
	"fmt"
	"sync"
)

// provenanceKey builds a unique registry key from mount and path.
func provenanceKey(mount, path string) string {
	return mount + "/" + path
}

// SecretProvenanceRegistry stores provenance records for secrets.
type SecretProvenanceRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretProvenance
}

// NewSecretProvenanceRegistry returns an initialised SecretProvenanceRegistry.
func NewSecretProvenanceRegistry() *SecretProvenanceRegistry {
	return &SecretProvenanceRegistry{
		entries: make(map[string]SecretProvenance),
	}
}

// Set validates and stores a SecretProvenance record.
func (r *SecretProvenanceRegistry) Set(p SecretProvenance) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("provenance registry: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[provenanceKey(p.Mount, p.Path)] = p
	return nil
}

// Get retrieves a SecretProvenance record by mount and path.
func (r *SecretProvenanceRegistry) Get(mount, path string) (SecretProvenance, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.entries[provenanceKey(mount, path)]
	return p, ok
}

// Remove deletes a provenance record from the registry.
func (r *SecretProvenanceRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, provenanceKey(mount, path))
}

// All returns a snapshot of all stored provenance records.
func (r *SecretProvenanceRegistry) All() []SecretProvenance {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretProvenance, 0, len(r.entries))
	for _, p := range r.entries {
		out = append(out, p)
	}
	return out
}
