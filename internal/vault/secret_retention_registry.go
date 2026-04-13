package vault

import (
	"fmt"
	"sync"
	"time"
)

func retentionKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretRetentionRegistry stores retention policies keyed by mount+path.
type SecretRetentionRegistry struct {
	mu       sync.RWMutex
	policies map[string]RetentionPolicy
}

// NewSecretRetentionRegistry returns an initialised registry.
func NewSecretRetentionRegistry() *SecretRetentionRegistry {
	return &SecretRetentionRegistry{
		policies: make(map[string]RetentionPolicy),
	}
}

// Set validates and stores a retention policy.
func (r *SecretRetentionRegistry) Set(p RetentionPolicy) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.policies[retentionKey(p.Mount, p.Path)] = p
	return nil
}

// Get retrieves a retention policy by mount and path.
func (r *SecretRetentionRegistry) Get(mount, path string) (RetentionPolicy, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.policies[retentionKey(mount, path)]
	return p, ok
}

// Remove deletes a retention policy from the registry.
func (r *SecretRetentionRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.policies, retentionKey(mount, path))
}

// Expired returns all policies whose retention window has passed.
func (r *SecretRetentionRegistry) Expired() []RetentionPolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []RetentionPolicy
	for _, p := range r.policies {
		if p.IsExpired() {
			out = append(out, p)
		}
	}
	return out
}

// All returns every registered retention policy.
func (r *SecretRetentionRegistry) All() []RetentionPolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]RetentionPolicy, 0, len(r.policies))
	for _, p := range r.policies {
		out = append(out, p)
	}
	return out
}
