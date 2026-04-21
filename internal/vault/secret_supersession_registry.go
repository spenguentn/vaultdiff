package vault

import (
	"fmt"
	"sync"
	"time"
)

func supersessionKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretSupersessionRegistry tracks secrets that have been superseded by newer versions or paths.
type SecretSupersessionRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretSupersession
}

// NewSecretSupersessionRegistry returns an initialised SecretSupersessionRegistry.
func NewSecretSupersessionRegistry() *SecretSupersessionRegistry {
	return &SecretSupersessionRegistry{
		entries: make(map[string]SecretSupersession),
	}
}

// Record marks a secret as superseded. SupersededAt is set automatically if zero.
func (r *SecretSupersessionRegistry) Record(s SecretSupersession) error {
	if err := s.Validate(); err != nil {
		return err
	}
	if s.SupersededAt.IsZero() {
		s.SupersededAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[supersessionKey(s.Mount, s.Path)] = s
	return nil
}

// Get returns the supersession record for the given mount and path.
func (r *SecretSupersessionRegistry) Get(mount, path string) (SecretSupersession, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.entries[supersessionKey(mount, path)]
	return s, ok
}

// Remove deletes the supersession record for the given mount and path.
func (r *SecretSupersessionRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, supersessionKey(mount, path))
}

// All returns a snapshot of all recorded supersession entries.
func (r *SecretSupersessionRegistry) All() []SecretSupersession {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretSupersession, 0, len(r.entries))
	for _, s := range r.entries {
		out = append(out, s)
	}
	return out
}
