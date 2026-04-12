package vault

import (
	"fmt"
	"sync"
	"time"
)

func auditTrailKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretAuditTrailRegistry stores ordered audit trail entries per secret.
type SecretAuditTrailRegistry struct {
	mu      sync.RWMutex
	entries map[string][]AuditTrailEntry
}

// NewSecretAuditTrailRegistry returns an initialised registry.
func NewSecretAuditTrailRegistry() *SecretAuditTrailRegistry {
	return &SecretAuditTrailRegistry{
		entries: make(map[string][]AuditTrailEntry),
	}
}

// Record appends a new entry for the secret identified by mount+path.
// The entry's Timestamp is set to now if it is zero.
func (r *SecretAuditTrailRegistry) Record(entry AuditTrailEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := auditTrailKey(entry.Mount, entry.Path)
	r.entries[key] = append(r.entries[key], entry)
	return nil
}

// Get returns all audit trail entries for the given mount+path, oldest first.
func (r *SecretAuditTrailRegistry) Get(mount, path string) []AuditTrailEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := auditTrailKey(mount, path)
	out := make([]AuditTrailEntry, len(r.entries[key]))
	copy(out, r.entries[key])
	return out
}

// Clear removes all audit entries for the given mount+path.
func (r *SecretAuditTrailRegistry) Clear(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, auditTrailKey(mount, path))
}

// Len returns the total number of entries recorded across all secrets.
func (r *SecretAuditTrailRegistry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := 0
	for _, v := range r.entries {
		total += len(v)
	}
	return total
}
