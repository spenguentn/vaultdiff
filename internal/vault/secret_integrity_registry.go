package vault

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func integrityKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretIntegrityRegistry stores and retrieves integrity records.
type SecretIntegrityRegistry struct {
	mu      sync.RWMutex
	records map[string]*SecretIntegrityRecord
}

// NewSecretIntegrityRegistry returns an initialised registry.
func NewSecretIntegrityRegistry() *SecretIntegrityRegistry {
	return &SecretIntegrityRegistry{
		records: make(map[string]*SecretIntegrityRecord),
	}
}

// Record stores an integrity record, stamping CheckedAt when zero.
func (r *SecretIntegrityRegistry) Record(rec *SecretIntegrityRecord) error {
	if rec == nil {
		return errors.New("integrity registry: record must not be nil")
	}
	if err := rec.Validate(); err != nil {
		return err
	}
	if rec.CheckedAt.IsZero() {
		rec.CheckedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[integrityKey(rec.Mount, rec.Path)] = rec
	return nil
}

// Get returns the integrity record for the given mount and path.
func (r *SecretIntegrityRegistry) Get(mount, path string) (*SecretIntegrityRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.records[integrityKey(mount, path)]
	return rec, ok
}

// Remove deletes the integrity record for the given mount and path.
func (r *SecretIntegrityRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, integrityKey(mount, path))
}

// All returns a snapshot of every stored record.
func (r *SecretIntegrityRegistry) All() []*SecretIntegrityRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*SecretIntegrityRecord, 0, len(r.records))
	for _, rec := range r.records {
		out = append(out, rec)
	}
	return out
}
