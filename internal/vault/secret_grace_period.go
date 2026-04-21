package vault

import (
	"errors"
	"fmt"
	"time"
)

// GracePeriodStatus represents the current state of a secret's grace period.
type GracePeriodStatus string

const (
	GracePeriodActive  GracePeriodStatus = "active"
	GracePeriodExpired GracePeriodStatus = "expired"
	GracePeriodPending GracePeriodStatus = "pending"
)

// IsValidGracePeriodStatus returns true if the provided status is a known value.
func IsValidGracePeriodStatus(s GracePeriodStatus) bool {
	switch s {
	case GracePeriodActive, GracePeriodExpired, GracePeriodPending:
		return true
	}
	return false
}

// SecretGracePeriod tracks a grace window during which a secret remains
// accessible after its scheduled expiry or rotation deadline.
type SecretGracePeriod struct {
	Mount      string            `json:"mount"`
	Path       string            `json:"path"`
	Status     GracePeriodStatus `json:"status"`
	StartsAt   time.Time         `json:"starts_at"`
	Duration   time.Duration     `json:"duration"`
	GrantedBy  string            `json:"granted_by"`
	Reason     string            `json:"reason,omitempty"`
}

// FullPath returns the canonical mount+path identifier.
func (g SecretGracePeriod) FullPath() string {
	return fmt.Sprintf("%s/%s", g.Mount, g.Path)
}

// ExpiresAt returns the time at which the grace period ends.
func (g SecretGracePeriod) ExpiresAt() time.Time {
	return g.StartsAt.Add(g.Duration)
}

// IsExpired reports whether the grace period has passed.
func (g SecretGracePeriod) IsExpired() bool {
	return time.Now().After(g.ExpiresAt())
}

// Validate checks that required fields are present and valid.
func (g SecretGracePeriod) Validate() error {
	if g.Mount == "" {
		return errors.New("grace period: mount is required")
	}
	if g.Path == "" {
		return errors.New("grace period: path is required")
	}
	if !IsValidGracePeriodStatus(g.Status) {
		return fmt.Errorf("grace period: unknown status %q", g.Status)
	}
	if g.Duration <= 0 {
		return errors.New("grace period: duration must be positive")
	}
	if g.GrantedBy == "" {
		return errors.New("grace period: granted_by is required")
	}
	return nil
}
