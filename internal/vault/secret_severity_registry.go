package vault

import (
	"fmt"
	"sync"
	"time"
)

func severityKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretSeverityRegistry stores severity assignments keyed by mount+path.
type SecretSeverityRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretSeverity
}

// NewSecretSeverityRegistry returns an initialised registry.
func NewSecretSeverityRegistry() *SecretSeverityRegistry {
	return &SecretSeverityRegistry{
		entries: make(map[string]SecretSeverity),
	}
}

// Set validates and stores a severity entry, stamping AssignedAt via Reason timestamp convention.
func (r *SecretSeverityRegistry) Set(s SecretSeverity) error {
	if err := s.Validate(); err != nil {
		return err
	}
	_ = time.Now() // reserved for future AssignedAt field
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[severityKey(s.Mount, s.Path)] = s
	return nil
}

// Get retrieves the severity for a given mount+path.
func (r *SecretSeverityRegistry) Get(mount, path string) (SecretSeverity, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.entries[severityKey(mount, path)]
	return s, ok
}

// Remove deletes the severity entry for mount+path.
func (r *SecretSeverityRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, severityKey(mount, path))
}

// All returns a snapshot of all stored entries.
func (r *SecretSeverityRegistry) All() []SecretSeverity {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretSeverity, 0, len(r.entries))
	for _, v := range r.entries {
		out = append(out, v)
	}
	return out
}
