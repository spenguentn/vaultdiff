package vault

import (
	"fmt"
	"sync"
	"time"
)

func consentKey(mount, path, grantedTo string) string {
	return fmt.Sprintf("%s/%s::%s", mount, path, grantedTo)
}

// SecretConsentRegistry stores and retrieves SecretConsent records.
type SecretConsentRegistry struct {
	mu      sync.RWMutex
	records map[string]SecretConsent
}

// NewSecretConsentRegistry returns an initialised SecretConsentRegistry.
func NewSecretConsentRegistry() *SecretConsentRegistry {
	return &SecretConsentRegistry{
		records: make(map[string]SecretConsent),
	}
}

// Grant stores a new consent record, stamping GrantedAt if zero.
func (r *SecretConsentRegistry) Grant(c SecretConsent) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.GrantedAt.IsZero() {
		c.GrantedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[consentKey(c.Mount, c.Path, c.GrantedTo)] = c
	return nil
}

// Get retrieves a consent record for the given mount, path and grantee.
func (r *SecretConsentRegistry) Get(mount, path, grantedTo string) (SecretConsent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.records[consentKey(mount, path, grantedTo)]
	return c, ok
}

// Revoke marks an existing consent record as revoked.
func (r *SecretConsentRegistry) Revoke(mount, path, grantedTo string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := consentKey(mount, path, grantedTo)
	c, ok := r.records[key]
	if !ok {
		return fmt.Errorf("consent: record not found for %s", key)
	}
	c.Status = ConsentRevoked
	r.records[key] = c
	return nil
}

// Remove deletes a consent record entirely.
func (r *SecretConsentRegistry) Remove(mount, path, grantedTo string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, consentKey(mount, path, grantedTo))
}

// All returns a snapshot of every stored consent record.
func (r *SecretConsentRegistry) All() []SecretConsent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretConsent, 0, len(r.records))
	for _, c := range r.records {
		out = append(out, c)
	}
	return out
}
