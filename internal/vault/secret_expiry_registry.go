package vault

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func expiryKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretExpiryRegistry tracks expiry policies for secrets.
type SecretExpiryRegistry struct {
	mu       sync.RWMutex
	policies map[string]ExpiryPolicy
}

// NewSecretExpiryRegistry returns an initialised registry.
func NewSecretExpiryRegistry() *SecretExpiryRegistry {
	return &SecretExpiryRegistry{
		policies: make(map[string]ExpiryPolicy),
	}
}

// Register adds or replaces an expiry policy.
func (r *SecretExpiryRegistry) Register(p ExpiryPolicy) error {
	if err := p.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.policies[expiryKey(p.Mount, p.Path)] = p
	return nil
}

// Get returns the policy for the given mount and path.
func (r *SecretExpiryRegistry) Get(mount, path string) (ExpiryPolicy, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.policies[expiryKey(mount, path)]
	return p, ok
}

// Remove deletes a policy from the registry.
func (r *SecretExpiryRegistry) Remove(mount, path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := expiryKey(mount, path)
	if _, ok := r.policies[k]; !ok {
		return errors.New("expiry registry: policy not found")
	}
	delete(r.policies, k)
	return nil
}

// CheckAll evaluates all registered policies against now and returns statuses.
func (r *SecretExpiryRegistry) CheckAll(now time.Time) []ExpiryStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]ExpiryStatus, 0, len(r.policies))
	for _, p := range r.policies {
		out = append(out, ExpiryStatus{
			Policy:  p,
			Expired: p.IsExpired(now),
			Soon:    p.IsExpiringSoon(now),
			Checked: now,
		})
	}
	return out
}

// Count returns the number of registered policies.
func (r *SecretExpiryRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.policies)
}
