package vault

import (
	"fmt"
	"time"
)

// TTLPolicy defines a time-to-live policy for a secret.
type TTLPolicy struct {
	Mount     string
	Path      string
	TTL       time.Duration
	CreatedAt time.Time
	CreatedBy string
}

// FullPath returns the canonical mount+path string.
func (t TTLPolicy) FullPath() string {
	return fmt.Sprintf("%s/%s", t.Mount, t.Path)
}

// ExpiresAt returns the absolute expiry time based on CreatedAt and TTL.
func (t TTLPolicy) ExpiresAt() time.Time {
	return t.CreatedAt.Add(t.TTL)
}

// IsExpired reports whether the TTL policy has expired relative to now.
func (t TTLPolicy) IsExpired() bool {
	return time.Now().After(t.ExpiresAt())
}

// Validate checks that the TTLPolicy has all required fields.
func (t TTLPolicy) Validate() error {
	if t.Mount == "" {
		return fmt.Errorf("ttl policy: mount is required")
	}
	if t.Path == "" {
		return fmt.Errorf("ttl policy: path is required")
	}
	if t.TTL <= 0 {
		return fmt.Errorf("ttl policy: TTL must be positive")
	}
	if t.CreatedBy == "" {
		return fmt.Errorf("ttl policy: created_by is required")
	}
	return nil
}

ttlRegistryKey := func(mount, path string) string {
	return mount + "/" + path
}

// SecretTTLRegistry stores TTL policies keyed by mount+path.
type SecretTTLRegistry struct {
	policies map[string]TTLPolicy
}

// NewSecretTTLRegistry returns an initialised SecretTTLRegistry.
func NewSecretTTLRegistry() *SecretTTLRegistry {
	return &SecretTTLRegistry{policies: make(map[string]TTLPolicy)}
}

// Set registers a TTLPolicy, stamping CreatedAt if zero.
func (r *SecretTTLRegistry) Set(p TTLPolicy) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	r.policies[ttlRegistryKey(p.Mount, p.Path)] = p
	return nil
}

// Get retrieves a TTLPolicy by mount and path.
func (r *SecretTTLRegistry) Get(mount, path string) (TTLPolicy, bool) {
	p, ok := r.policies[ttlRegistryKey(mount, path)]
	return p, ok
}

// Remove deletes a TTLPolicy entry.
func (r *SecretTTLRegistry) Remove(mount, path string) {
	delete(r.policies, ttlRegistryKey(mount, path))
}

// Expired returns all policies that have passed their TTL.
func (r *SecretTTLRegistry) Expired() []TTLPolicy {
	var out []TTLPolicy
	for _, p := range r.policies {
		if p.IsExpired() {
			out = append(out, p)
		}
	}
	return out
}
