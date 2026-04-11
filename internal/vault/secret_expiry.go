package vault

import (
	"errors"
	"fmt"
	"time"
)

// ExpiryPolicy defines when a secret is considered expired or expiring soon.
type ExpiryPolicy struct {
	Mount       string
	Path        string
	ExpiresAt   time.Time
	WarnBefore  time.Duration
	Owner       string
}

// Validate checks that the expiry policy is well-formed.
func (e ExpiryPolicy) Validate() error {
	if e.Mount == "" {
		return errors.New("expiry policy: mount is required")
	}
	if e.Path == "" {
		return errors.New("expiry policy: path is required")
	}
	if e.ExpiresAt.IsZero() {
		return errors.New("expiry policy: expires_at must be set")
	}
	if e.WarnBefore < 0 {
		return errors.New("expiry policy: warn_before must not be negative")
	}
	return nil
}

// FullPath returns the canonical mount+path string.
func (e ExpiryPolicy) FullPath() string {
	return fmt.Sprintf("%s/%s", e.Mount, e.Path)
}

// IsExpired reports whether the secret has passed its expiry time.
func (e ExpiryPolicy) IsExpired(now time.Time) bool {
	return !e.ExpiresAt.IsZero() && now.After(e.ExpiresAt)
}

// IsExpiringSoon reports whether the secret will expire within the warn window.
func (e ExpiryPolicy) IsExpiringSoon(now time.Time) bool {
	if e.ExpiresAt.IsZero() || e.WarnBefore == 0 {
		return false
	}
	return now.Add(e.WarnBefore).After(e.ExpiresAt)
}

// ExpiryStatus summarises the expiry state of a secret.
type ExpiryStatus struct {
	Policy  ExpiryPolicy
	Expired bool
	Soon    bool
	Checked time.Time
}

// String returns a human-readable status label.
func (s ExpiryStatus) String() string {
	switch {
	case s.Expired:
		return "expired"
	case s.Soon:
		return "expiring-soon"
	default:
		return "ok"
	}
}
