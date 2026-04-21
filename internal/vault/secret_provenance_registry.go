package vault

import (
	"fmt"
	"sync"
	"time"
)

func provenanceKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretProvenanceRegistry stores provenance records keyed by mount+path.
type SecretProvenanceRegistry struct {
	mu      sync.RWMutex
	records map[string]SecretProvenance
}

// NewSecretProvenanceRegistry returns an initialised registry.
func NewSecretProvenanceRegistry() *SecretProvenanceRegistry {
	return &SecretProvenanceRegistry{
		records: make(map[string]SecretProvenance),
	}
}

// Set stores a provenance record, stamping RecordedAt if zero.
func (r *SecretProvenanceRegistry) Set(p SecretProvenance) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if p.RecordedAt.IsZero() {
		p.RecordedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[provenanceKey(p.Mount, p.Path)] = p
	return nil
}

// Get retrieves the provenance record for a given mount+path.
func (r *SecretProvenanceRegistry) Get(mount, path string) (SecretProvenance, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.records[provenanceKey(mount, path)]
	return p, ok
}

// Remove deletes the provenance record for a given mount+path.
func (r *SecretProvenanceRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, provenanceKey(mount, path))
}

// All returns a snapshot of all stored provenance records.
func (r *SecretProvenanceRegistry) All() []SecretProvenance {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretProvenance, 0, len(r.records))
	for _, v := range r.records {
		out = append(out, v)
	}
	return out
}
