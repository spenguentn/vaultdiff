package vault

import (
	"fmt"
	"sync"
	"time"
)

func delegationKey(mount, path, delegateTo string) string {
	return fmt.Sprintf("%s/%s::%s", mount, path, delegateTo)
}

// SecretDelegationRegistry manages active delegations for secrets.
type SecretDelegationRegistry struct {
	mu          sync.RWMutex
	delegations map[string]SecretDelegation
}

// NewSecretDelegationRegistry returns an initialised registry.
func NewSecretDelegationRegistry() *SecretDelegationRegistry {
	return &SecretDelegationRegistry{
		delegations: make(map[string]SecretDelegation),
	}
}

// Delegate stores a delegation, stamping DelegatedAt if zero.
func (r *SecretDelegationRegistry) Delegate(d SecretDelegation) error {
	if err := d.Validate(); err != nil {
		return err
	}
	if d.DelegatedAt.IsZero() {
		d.DelegatedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.delegations[delegationKey(d.Mount, d.Path, d.DelegateTo)] = d
	return nil
}

// Get retrieves a delegation by mount, path and delegateTo.
func (r *SecretDelegationRegistry) Get(mount, path, delegateTo string) (SecretDelegation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.delegations[delegationKey(mount, path, delegateTo)]
	return d, ok
}

// Revoke removes a delegation entry.
func (r *SecretDelegationRegistry) Revoke(mount, path, delegateTo string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.delegations, delegationKey(mount, path, delegateTo))
}

// Active returns all non-expired delegations.
func (r *SecretDelegationRegistry) Active() []SecretDelegation {
	r.mu.RLock()
	defer r.mu.RUnlock()
	now := time.Now().UTC()
	out := make([]SecretDelegation, 0)
	for _, d := range r.delegations {
		if d.ExpiresAt.IsZero() || d.ExpiresAt.After(now) {
			out = append(out, d)
		}
	}
	return out
}
