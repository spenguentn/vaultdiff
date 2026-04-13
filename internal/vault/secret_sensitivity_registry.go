package vault

import (
	"fmt"
	"sync"
)

func sensitivityKey(mount, path string) string {
	return mount + "::" + path
}

// SecretSensitivityRegistry stores sensitivity records for secrets.
type SecretSensitivityRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretSensitivity
}

// NewSecretSensitivityRegistry returns an initialised registry.
func NewSecretSensitivityRegistry() *SecretSensitivityRegistry {
	return &SecretSensitivityRegistry{
		entries: make(map[string]SecretSensitivity),
	}
}

// Set stores or updates the sensitivity record after validation.
func (r *SecretSensitivityRegistry) Set(s SecretSensitivity) error {
	if err := s.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[sensitivityKey(s.Mount, s.Path)] = s
	return nil
}

// Get retrieves the sensitivity record for a given mount and path.
func (r *SecretSensitivityRegistry) Get(mount, path string) (SecretSensitivity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.entries[sensitivityKey(mount, path)]
	if !ok {
		return SecretSensitivity{}, fmt.Errorf("sensitivity: no record for %s/%s", mount, path)
	}
	return s, nil
}

// Remove deletes the sensitivity record for a given mount and path.
func (r *SecretSensitivityRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, sensitivityKey(mount, path))
}

// All returns a snapshot of every stored sensitivity record.
func (r *SecretSensitivityRegistry) All() []SecretSensitivity {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretSensitivity, 0, len(r.entries))
	for _, s := range r.entries {
		out = append(out, s)
	}
	return out
}
