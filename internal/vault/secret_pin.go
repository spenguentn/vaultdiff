package vault

import (
	"fmt"
	"time"
)

// SecretPin represents a pinned version of a secret, preventing it from
// being rotated or overwritten until explicitly unpinned.
type SecretPin struct {
	Mount     string    `json:"mount"`
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	PinnedBy  string    `json:"pinned_by"`
	Reason    string    `json:"reason"`
	PinnedAt  time.Time `json:"pinned_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// FullPath returns the canonical mount+path string.
func (p SecretPin) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Mount, p.Path)
}

// IsExpired reports whether the pin has passed its expiry time.
func (p SecretPin) IsExpired() bool {
	if p.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(p.ExpiresAt)
}

// Validate checks that required fields are present.
func (p SecretPin) Validate() error {
	if p.Mount == "" {
		return fmt.Errorf("secret pin: mount is required")
	}
	if p.Path == "" {
		return fmt.Errorf("secret pin: path is required")
	}
	if p.Version <= 0 {
		return fmt.Errorf("secret pin: version must be greater than zero")
	}
	if p.PinnedBy == "" {
		return fmt.Errorf("secret pin: pinned_by is required")
	}
	if p.Reason == "" {
		return fmt.Errorf("secret pin: reason is required")
	}
	return nil
}
