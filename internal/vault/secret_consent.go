package vault

import (
	"errors"
	"fmt"
	"time"
)

// ConsentStatus represents the approval state of a secret access consent.
type ConsentStatus string

const (
	ConsentPending  ConsentStatus = "pending"
	ConsentGranted  ConsentStatus = "granted"
	ConsentRevoked  ConsentStatus = "revoked"
	ConsentExpired  ConsentStatus = "expired"
)

// IsValidConsentStatus reports whether s is a recognised ConsentStatus.
func IsValidConsentStatus(s ConsentStatus) bool {
	switch s {
	case ConsentPending, ConsentGranted, ConsentRevoked, ConsentExpired:
		return true
	}
	return false
}

// SecretConsent records an explicit consent decision for accessing a secret.
type SecretConsent struct {
	Mount      string        `json:"mount"`
	Path       string        `json:"path"`
	GrantedTo  string        `json:"granted_to"`
	GrantedBy  string        `json:"granted_by"`
	Status     ConsentStatus `json:"status"`
	GrantedAt  time.Time     `json:"granted_at"`
	ExpiresAt  *time.Time    `json:"expires_at,omitempty"`
	Reason     string        `json:"reason,omitempty"`
}

// FullPath returns the canonical vault path for this consent record.
func (c SecretConsent) FullPath() string {
	return fmt.Sprintf("%s/%s", c.Mount, c.Path)
}

// IsExpired reports whether the consent window has passed.
func (c SecretConsent) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// IsActive reports whether the consent is currently granted and not expired.
func (c SecretConsent) IsActive() bool {
	return c.Status == ConsentGranted && !c.IsExpired()
}

// Validate checks that the consent record is well-formed.
func (c SecretConsent) Validate() error {
	if c.Mount == "" {
		return errors.New("consent: mount is required")
	}
	if c.Path == "" {
		return errors.New("consent: path is required")
	}
	if c.GrantedTo == "" {
		return errors.New("consent: granted_to is required")
	}
	if c.GrantedBy == "" {
		return errors.New("consent: granted_by is required")
	}
	if !IsValidConsentStatus(c.Status) {
		return fmt.Errorf("consent: unknown status %q", c.Status)
	}
	if c.GrantedAt.IsZero() {
		return errors.New("consent: granted_at is required")
	}
	return nil
}
