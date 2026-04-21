package vault

import (
	"fmt"
	"sync"
	"time"
)

func traceabilityKey(mount, path string, version int) string {
	return fmt.Sprintf("%s/%s@v%d", mount, path, version)
}

// SecretTraceabilityRegistry stores and retrieves traceability records.
type SecretTraceabilityRegistry struct {
	mu      sync.RWMutex
	records map[string]*SecretTraceability
}

// NewSecretTraceabilityRegistry returns an initialised registry.
func NewSecretTraceabilityRegistry() *SecretTraceabilityRegistry {
	return &SecretTraceabilityRegistry{
		records: make(map[string]*SecretTraceability),
	}
}

// Record stores a traceability entry, setting TracedAt if zero.
func (r *SecretTraceabilityRegistry) Record(t *SecretTraceability) error {
	if err := t.Validate(); err != nil {
		return err
	}
	if t.TracedAt.IsZero() {
		t.TracedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[traceabilityKey(t.Mount, t.Path, t.Version)] = t
	return nil
}

// Get retrieves a traceability record by mount, path and version.
func (r *SecretTraceabilityRegistry) Get(mount, path string, version int) (*SecretTraceability, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.records[traceabilityKey(mount, path, version)]
	return v, ok
}

// Remove deletes a traceability record.
func (r *SecretTraceabilityRegistry) Remove(mount, path string, version int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, traceabilityKey(mount, path, version))
}

// Len returns the number of stored records.
func (r *SecretTraceabilityRegistry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.records)
}
