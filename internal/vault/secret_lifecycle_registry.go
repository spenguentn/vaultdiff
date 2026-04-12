package vault

import (
	"fmt"
	"sync"
	"time"
)

func lifecycleKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretLifecycleRegistry stores and retrieves lifecycle records for secrets.
type SecretLifecycleRegistry struct {
	mu      sync.RWMutex
	entries map[string]*SecretLifecycle
}

// NewSecretLifecycleRegistry returns an initialised registry.
func NewSecretLifecycleRegistry() *SecretLifecycleRegistry {
	return &SecretLifecycleRegistry{
		entries: make(map[string]*SecretLifecycle),
	}
}

// Set registers or replaces the lifecycle record for a secret.
// CreatedAt is set automatically if it is zero.
func (r *SecretLifecycleRegistry) Set(lc *SecretLifecycle) error {
	if err := lc.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if lc.CreatedAt.IsZero() {
		lc.CreatedAt = time.Now().UTC()
	}
	r.entries[lifecycleKey(lc.Mount, lc.Path)] = lc
	return nil
}

// Get retrieves the lifecycle record for the given mount and path.
func (r *SecretLifecycleRegistry) Get(mount, path string) (*SecretLifecycle, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.entries[lifecycleKey(mount, path)]
	return v, ok
}

// Remove deletes the lifecycle record for the given mount and path.
func (r *SecretLifecycleRegistry) Remove(mount, path string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := lifecycleKey(mount, path)
	_, ok := r.entries[key]
	if ok {
		delete(r.entries, key)
	}
	return ok
}

// All returns a snapshot of every registered lifecycle record.
func (r *SecretLifecycleRegistry) All() []*SecretLifecycle {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*SecretLifecycle, 0, len(r.entries))
	for _, v := range r.entries {
		out = append(out, v)
	}
	return out
}
