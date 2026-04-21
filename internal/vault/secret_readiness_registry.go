package vault

import (
	"fmt"
	"sync"
	"time"
)

func readinessKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretReadinessRegistry tracks readiness records for secrets.
type SecretReadinessRegistry struct {
	mu      sync.RWMutex
	records map[string]SecretReadiness
}

// NewSecretReadinessRegistry returns an initialised SecretReadinessRegistry.
func NewSecretReadinessRegistry() *SecretReadinessRegistry {
	return &SecretReadinessRegistry{
		records: make(map[string]SecretReadiness),
	}
}

// Set stores a readiness record, stamping AssessedAt if zero.
func (r *SecretReadinessRegistry) Set(rec SecretReadiness) error {
	if err := rec.Validate(); err != nil {
		return err
	}
	if rec.AssessedAt.IsZero() {
		rec.AssessedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[readinessKey(rec.Mount, rec.Path)] = rec
	return nil
}

// Get retrieves the readiness record for the given mount and path.
func (r *SecretReadinessRegistry) Get(mount, path string) (SecretReadiness, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.records[readinessKey(mount, path)]
	return rec, ok
}

// Remove deletes the readiness record for the given mount and path.
func (r *SecretReadinessRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, readinessKey(mount, path))
}

// All returns a snapshot of all stored readiness records.
func (r *SecretReadinessRegistry) All() []SecretReadiness {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretReadiness, 0, len(r.records))
	for _, v := range r.records {
		out = append(out, v)
	}
	return out
}
