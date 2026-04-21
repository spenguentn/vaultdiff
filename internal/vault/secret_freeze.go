package vault

import (
	"errors"
	"fmt"
	"time"
)

// FreezeReason describes why a secret was frozen.
type FreezeReason string

const (
	FreezeReasonManual     FreezeReason = "manual"
	FreezeReasonCompliance FreezeReason = "compliance"
	FreezeReasonIncident   FreezeReason = "incident"
)

// IsValidFreezeReason returns true if the reason is a known freeze reason.
func IsValidFreezeReason(r FreezeReason) bool {
	switch r {
	case FreezeReasonManual, FreezeReasonCompliance, FreezeReasonIncident:
		return true
	}
	return false
}

// SecretFreeze represents a freeze record applied to a secret path.
type SecretFreeze struct {
	Mount     string       `json:"mount"`
	Path      string       `json:"path"`
	FrozenBy  string       `json:"frozen_by"`
	Reason    FreezeReason `json:"reason"`
	Note      string       `json:"note,omitempty"`
	FrozenAt  time.Time    `json:"frozen_at"`
	ExpiresAt *time.Time   `json:"expires_at,omitempty"`
}

// FullPath returns the canonical mount+path string.
func (f SecretFreeze) FullPath() string {
	return fmt.Sprintf("%s/%s", f.Mount, f.Path)
}

// IsExpired returns true when an expiry is set and has passed.
func (f SecretFreeze) IsExpired() bool {
	if f.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*f.ExpiresAt)
}

// IsActive returns true when the freeze record is valid and has not expired.
func (f SecretFreeze) IsActive() bool {
	return f.Validate() == nil && !f.IsExpired()
}

// Validate checks that the freeze record contains required fields.
func (f SecretFreeze) Validate() error {
	if f.Mount == "" {
		return errors.New("freeze: mount is required")
	}
	if f.Path == "" {
		return errors.New("freeze: path is required")
	}
	if f.FrozenBy == "" {
		return errors.New("freeze: frozen_by is required")
	}
	if !IsValidFreezeReason(f.Reason) {
		return fmt.Errorf("freeze: unknown reason %q", f.Reason)
	}
	if f.ExpiresAt != nil && !f.ExpiresAt.After(f.FrozenAt) {
		return errors.New("freeze: expires_at must be after frozen_at")
	}
	return nil
}
