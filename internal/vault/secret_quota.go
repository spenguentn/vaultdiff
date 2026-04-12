package vault

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// QuotaScope defines the scope at which a quota applies.
type QuotaScope string

const (
	QuotaScopeMount  QuotaScope = "mount"
	QuotaScopePrefix QuotaScope = "prefix"
	QuotaScopeGlobal QuotaScope = "global"
)

// SecretQuota defines a read-rate or access quota for a secret path.
type SecretQuota struct {
	Mount      string        `json:"mount"`
	Prefix     string        `json:"prefix"`
	Scope      QuotaScope    `json:"scope"`
	MaxReads   int           `json:"max_reads"`
	WindowSize time.Duration `json:"window_size"`
	CreatedAt  time.Time     `json:"created_at"`
}

// FullPath returns a composite key for the quota.
func (q SecretQuota) FullPath() string {
	return fmt.Sprintf("%s/%s", q.Mount, q.Prefix)
}

// Validate checks that the quota is well-formed.
func (q SecretQuota) Validate() error {
	if q.Mount == "" {
		return errors.New("quota: mount is required")
	}
	if q.MaxReads <= 0 {
		return errors.New("quota: max_reads must be greater than zero")
	}
	if q.WindowSize <= 0 {
		return errors.New("quota: window_size must be greater than zero")
	}
	if q.Scope == "" {
		return errors.New("quota: scope is required")
	}
	return nil
}

// quotaKey builds a registry key.
func quotaKey(mount, prefix string) string {
	return fmt.Sprintf("%s::%s", mount, prefix)
}

// SecretQuotaRegistry stores and enforces secret quotas.
type SecretQuotaRegistry struct {
	mu     sync.RWMutex
	quotas map[string]SecretQuota
	counts map[string][]time.Time
}

// NewSecretQuotaRegistry returns an initialised registry.
func NewSecretQuotaRegistry() *SecretQuotaRegistry {
	return &SecretQuotaRegistry{
		quotas: make(map[string]SecretQuota),
		counts: make(map[string][]time.Time),
	}
}

// Register adds or replaces a quota.
func (r *SecretQuotaRegistry) Register(q SecretQuota) error {
	if err := q.Validate(); err != nil {
		return err
	}
	q.CreatedAt = time.Now().UTC()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.quotas[quotaKey(q.Mount, q.Prefix)] = q
	return nil
}

// Get retrieves a quota by mount and prefix.
func (r *SecretQuotaRegistry) Get(mount, prefix string) (SecretQuota, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	q, ok := r.quotas[quotaKey(mount, prefix)]
	return q, ok
}

// Remove deletes a quota entry.
func (r *SecretQuotaRegistry) Remove(mount, prefix string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.quotas, quotaKey(mount, prefix))
	delete(r.counts, quotaKey(mount, prefix))
}

// Allow records an access attempt and returns true if within quota.
func (r *SecretQuotaRegistry) Allow(mount, prefix string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := quotaKey(mount, prefix)
	q, ok := r.quotas[key]
	if !ok {
		return true
	}
	now := time.Now().UTC()
	cutoff := now.Add(-q.WindowSize)
	var valid []time.Time
	for _, t := range r.counts[key] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	if len(valid) >= q.MaxReads {
		r.counts[key] = valid
		return false
	}
	r.counts[key] = append(valid, now)
	return true
}
