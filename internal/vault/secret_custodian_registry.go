package vault

import (
	"fmt"
	"sync"
	"time"
)

func custodianKey(mount, path, custodian string) string {
	return fmt.Sprintf("%s/%s::%s", mount, path, custodian)
}

// SecretCustodianRegistry stores custodian assignments for secrets.
type SecretCustodianRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretCustodian
}

// NewSecretCustodianRegistry creates an empty custodian registry.
func NewSecretCustodianRegistry() *SecretCustodianRegistry {
	return &SecretCustodianRegistry{
		entries: make(map[string]SecretCustodian),
	}
}

// Assign adds or updates a custodian assignment. AssignedAt is set automatically.
func (r *SecretCustodianRegistry) Assign(c SecretCustodian) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.AssignedAt.IsZero() {
		c.AssignedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[custodianKey(c.Mount, c.Path, c.Custodian)] = c
	return nil
}

// Get retrieves a custodian assignment by mount, path, and custodian name.
func (r *SecretCustodianRegistry) Get(mount, path, custodian string) (SecretCustodian, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.entries[custodianKey(mount, path, custodian)]
	return c, ok
}

// Remove deletes a custodian assignment.
func (r *SecretCustodianRegistry) Remove(mount, path, custodian string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, custodianKey(mount, path, custodian))
}

// List returns all custodians assigned to a given mount+path.
func (r *SecretCustodianRegistry) List(mount, path string) []SecretCustodian {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []SecretCustodian
	for _, c := range r.entries {
		if c.Mount == mount && c.Path == path {
			result = append(result, c)
		}
	}
	return result
}
