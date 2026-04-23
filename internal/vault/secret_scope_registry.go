package vault

import (
	"fmt"
	"sync"
	"time"
)

func scopeKey(mount, path string) string {
	return mount + "/" + path
}

// SecretScopeRegistry stores and retrieves scope assignments for secrets.
type SecretScopeRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretScope
}

// NewSecretScopeRegistry returns an initialised SecretScopeRegistry.
func NewSecretScopeRegistry() *SecretScopeRegistry {
	return &SecretScopeRegistry{
		entries: make(map[string]SecretScope),
	}
}

// Set validates and stores a scope assignment, stamping AssignedAt if unset.
func (r *SecretScopeRegistry) Set(s SecretScope) error {
	if err := s.Validate(); err != nil {
		return fmt.Errorf("scope registry: %w", err)
	}
	if s.AssignedAt.IsZero() {
		s.AssignedAt = time.Now()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[scopeKey(s.Mount, s.Path)] = s
	return nil
}

// Get retrieves the scope assignment for the given mount and path.
func (r *SecretScopeRegistry) Get(mount, path string) (SecretScope, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.entries[scopeKey(mount, path)]
	return s, ok
}

// Remove deletes the scope assignment for the given mount and path.
func (r *SecretScopeRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, scopeKey(mount, path))
}
