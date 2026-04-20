package vault

import (
	"fmt"
	"time"
)

// CovenantType represents the type of agreement governing a secret.
type CovenantType string

const (
	CovenantTypeShared    CovenantType = "shared"
	CovenantTypeExclusive CovenantType = "exclusive"
	CovenantTypeReadOnly  CovenantType = "read_only"
)

// IsValidCovenantType returns true if the given type is a recognised covenant type.
func IsValidCovenantType(t CovenantType) bool {
	switch t {
	case CovenantTypeShared, CovenantTypeExclusive, CovenantTypeReadOnly:
		return true
	}
	return false
}

// SecretCovenant records a governance agreement attached to a secret.
type SecretCovenant struct {
	Mount       string       `json:"mount"`
	Path        string       `json:"path"`
	Type        CovenantType `json:"type"`
	Owner       string       `json:"owner"`
	Description string       `json:"description"`
	ExpiresAt   *time.Time   `json:"expires_at,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
}

// FullPath returns the canonical mount+path string for this covenant.
func (c *SecretCovenant) FullPath() string {
	return fmt.Sprintf("%s/%s", c.Mount, c.Path)
}

// IsExpired returns true when the covenant has a set expiry that is in the past.
func (c *SecretCovenant) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// Validate returns an error if the covenant is missing required fields or has
// an unrecognised type.
func (c *SecretCovenant) Validate() error {
	if c.Mount == "" {
		return fmt.Errorf("covenant: mount is required")
	}
	if c.Path == "" {
		return fmt.Errorf("covenant: path is required")
	}
	if c.Owner == "" {
		return fmt.Errorf("covenant: owner is required")
	}
	if !IsValidCovenantType(c.Type) {
		return fmt.Errorf("covenant: unknown type %q", c.Type)
	}
	return nil
}
