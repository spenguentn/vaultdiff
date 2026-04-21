package vault

import (
	"fmt"
	"sync"
	"time"
)

func revocationKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretRevocationRegistry stores and retrieves revocation records by mount+path.
type SecretRevocationRegistry struct {
	mu		sync.RWMutex
	records	map[string]RevocationRecord
}

// NewSecretRevocationRegistry returns an initialised SecretRevocationRegistry.
func NewSecretRevocationRegistry() *SecretRevocationRegistry {
	return &SecretRevocationRegistry{
		records: make(map[string]RevocationRecord),
	}
}

// Revoke stores a revocation record, setting RevokedAt if it is zero.
func (r *SecretRevocationRegistry) Revoke(rec RevocationRecord) error {
	if rec.RevokedAt.IsZero() {
		rec.RevokedAt = time.Now().UTC()
	}
	if err := rec.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[revocationKey(rec.Mount, rec.Path)] = rec
	return nil
}

// Get returns the revocation record for the given mount and path.
func (r *SecretRevocationRegistry) Get(mount, path string) (RevocationRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.records[revocationKey(mount, path)]
	return rec, ok
}

// Remove deletes the revocation record for the given mount and path.
func (r *SecretRevocationRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, revocationKey(mount, path))
}

// All returns a slice of all stored revocation records.
func (r *SecretRevocationRegistry) All() []RevocationRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]RevocationRecord, 0, len(r.records))
	for _, rec := range r.records {
		out = append(out, rec)
	}
	return out
}
