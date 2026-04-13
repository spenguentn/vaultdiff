package vault

import (
	"fmt"
	"sync"
	"time"
)

func trustKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretTrustRegistry stores and retrieves trust records for secrets.
type SecretTrustRegistry struct {
	mu      sync.RWMutex
	entries map[string]*SecretTrust
}

// NewSecretTrustRegistry creates an empty SecretTrustRegistry.
func NewSecretTrustRegistry() *SecretTrustRegistry {
	return &SecretTrustRegistry{
		entries: make(map[string]*SecretTrust),
	}
}

// Set validates and stores a trust record, setting AssignedAt if zero.
func (r *SecretTrustRegistry) Set(t *SecretTrust) error {
	if err := t.Validate(); err != nil {
		return err
	}
	if t.AssignedAt.IsZero() {
		t.AssignedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[trustKey(t.Mount, t.Path)] = t
	return nil
}

// Get retrieves a trust record by mount and path.
func (r *SecretTrustRegistry) Get(mount, path string) (*SecretTrust, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.entries[trustKey(mount, path)]
	return t, ok
}

// Remove deletes a trust record. Returns false if not found.
func (r *SecretTrustRegistry) Remove(mount, path string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := trustKey(mount, path)
	if _, ok := r.entries[key]; !ok {
		return false
	}
	delete(r.entries, key)
	return true
}

// All returns a snapshot of all stored trust records.
func (r *SecretTrustRegistry) All() []*SecretTrust {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*SecretTrust, 0, len(r.entries))
	for _, v := range r.entries {
		out = append(out, v)
	}
	return out
}
