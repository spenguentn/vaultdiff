package vault

import (
	"fmt"
	"sync"
	"time"
)

func freezeKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretFreezeRegistry stores and manages secret freeze records.
type SecretFreezeRegistry struct {
	mu      sync.RWMutex
	records map[string]SecretFreeze
}

// NewSecretFreezeRegistry returns an initialised SecretFreezeRegistry.
func NewSecretFreezeRegistry() *SecretFreezeRegistry {
	return &SecretFreezeRegistry{
		records: make(map[string]SecretFreeze),
	}
}

// Freeze adds or replaces a freeze record after validation.
func (r *SecretFreezeRegistry) Freeze(f SecretFreeze) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if f.FrozenAt.IsZero() {
		f.FrozenAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[freezeKey(f.Mount, f.Path)] = f
	return nil
}

// Get returns the freeze record for the given mount and path.
func (r *SecretFreezeRegistry) Get(mount, path string) (SecretFreeze, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.records[freezeKey(mount, path)]
	return f, ok
}

// IsFrozen returns true when an active (non-expired) freeze exists.
func (r *SecretFreezeRegistry) IsFrozen(mount, path string) bool {
	f, ok := r.Get(mount, path)
	return ok && !f.IsExpired()
}

// Unfreeze removes a freeze record for the given mount and path.
func (r *SecretFreezeRegistry) Unfreeze(mount, path string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := freezeKey(mount, path)
	_, ok := r.records[key]
	if ok {
		delete(r.records, key)
	}
	return ok
}

// All returns a snapshot of all current freeze records.
func (r *SecretFreezeRegistry) All() []SecretFreeze {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretFreeze, 0, len(r.records))
	for _, f := range r.records {
		out = append(out, f)
	}
	return out
}
