package vault

import (
	"fmt"
	"sync"
	"time"
)

type SecretScopeRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretScope
}

func scopeKey(mount, path string) string {
	return fmt.Sprintf("%s::%s", mount, path)
}

func NewSecretScopeRegistry() *SecretScopeRegistry {
	return &SecretScopeRegistry{
		entries: make(map[string]SecretScope),
	}
}

func (r *SecretScopeRegistry) Set(s SecretScope) error {
	if err := s.Validate(); err != nil {
		return err
	}
	if s.AssignedAt.IsZero() {
		s.AssignedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[scopeKey(s.Mount, s.Path)] = s
	return nil
}

func (r *SecretScopeRegistry) Get(mount, path string) (SecretScope, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.entries[scopeKey(mount, path)]
	return v, ok
}

func (r *SecretScopeRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, scopeKey(mount, path))
}

func (r *SecretScopeRegistry) All() []SecretScope {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretScope, 0, len(r.entries))
	for _, v := range r.entries {
		out = append(out, v)
	}
	return out
}
