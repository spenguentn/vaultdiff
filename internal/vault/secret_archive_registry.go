package vault

import (
	"fmt"
	"sync"
	"time"
)

// SecretArchiveRegistry stores and retrieves archived secret entries.
type SecretArchiveRegistry struct {
	mu      sync.RWMutex
	entries map[string]*SecretArchiveEntry
}

// NewSecretArchiveRegistry returns an initialised SecretArchiveRegistry.
func NewSecretArchiveRegistry() *SecretArchiveRegistry {
	return &SecretArchiveRegistry{
		entries: make(map[string]*SecretArchiveEntry),
	}
}

func archiveKey(mount, path string) string {
	return mount + ":" + path
}

// Archive validates and stores an entry. ArchivedAt is set if zero.
func (r *SecretArchiveRegistry) Archive(entry *SecretArchiveEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	if !IsReasonValid(entry.Reason) {
		return fmt.Errorf("archive registry: unknown reason %q", entry.Reason)
	}
	if entry.ArchivedAt.IsZero() {
		entry.ArchivedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[archiveKey(entry.Mount, entry.Path)] = entry
	return nil
}

// Get returns the archive entry for the given mount and path, if present.
func (r *SecretArchiveRegistry) Get(mount, path string) (*SecretArchiveEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[archiveKey(mount, path)]
	return e, ok
}

// Remove deletes an entry from the registry. Returns false if not found.
func (r *SecretArchiveRegistry) Remove(mount, path string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := archiveKey(mount, path)
	_, ok := r.entries[key]
	if ok {
		delete(r.entries, key)
	}
	return ok
}

// All returns a snapshot of all archived entries.
func (r *SecretArchiveRegistry) All() []*SecretArchiveEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*SecretArchiveEntry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e)
	}
	return out
}

// Count returns the number of archived entries.
func (r *SecretArchiveRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}
