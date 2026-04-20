package vault

import (
	"fmt"
	"sync"
	"time"
)

func obsolescenceKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretObsolescenceRegistry stores obsolescence records keyed by mount+path.
type SecretObsolescenceRegistry struct {
	mu      sync.RWMutex
	entries map[string]*SecretObsolescence
}

// NewSecretObsolescenceRegistry returns an initialised registry.
func NewSecretObsolescenceRegistry() *SecretObsolescenceRegistry {
	return &SecretObsolescenceRegistry{
		entries: make(map[string]*SecretObsolescence),
	}
}

// Mark records a secret as obsolete, setting MarkedAt if unset.
func (r *SecretObsolescenceRegistry) Mark(entry *SecretObsolescence) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	if entry.MarkedAt.IsZero() {
		entry.MarkedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[obsolescenceKey(entry.Mount, entry.Path)] = entry
	return nil
}

// Get retrieves the obsolescence record for the given mount+path.
func (r *SecretObsolescenceRegistry) Get(mount, path string) (*SecretObsolescence, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[obsolescenceKey(mount, path)]
	return e, ok
}

// Remove deletes the obsolescence record for the given mount+path.
func (r *SecretObsolescenceRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, obsolescenceKey(mount, path))
}

// All returns a snapshot of all registered obsolescence records.
func (r *SecretObsolescenceRegistry) All() []*SecretObsolescence {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*SecretObsolescence, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e)
	}
	return out
}
