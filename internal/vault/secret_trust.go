package vault

import (
	"fmt"
	"time"
)

// TrustLevel represents the degree of trust assigned to a secret.
type TrustLevel string

const (
	TrustLevelUntrusted TrustLevel = "untrusted"
	TrustLevelLow       TrustLevel = "low"
	TrustLevelMedium    TrustLevel = "medium"
	TrustLevelHigh      TrustLevel = "high"
	TrustLevelVerified  TrustLevel = "verified"
)

// IsValidTrustLevel returns true if the given level is a known trust level.
func IsValidTrustLevel(level TrustLevel) bool {
	switch level {
	case TrustLevelUntrusted, TrustLevelLow, TrustLevelMedium, TrustLevelHigh, TrustLevelVerified:
		return true
	}
	return false
}

// SecretTrust records the trust level assigned to a secret at a given path.
type SecretTrust struct {
	Mount      string     `json:"mount"`
	Path       string     `json:"path"`
	Level      TrustLevel `json:"level"`
	AssignedBy string     `json:"assigned_by"`
	AssignedAt time.Time  `json:"assigned_at"`
	Note       string     `json:"note,omitempty"`
}

// FullPath returns the canonical mount+path string.
func (t *SecretTrust) FullPath() string {
	return fmt.Sprintf("%s/%s", t.Mount, t.Path)
}

// Validate checks that the SecretTrust record is well-formed.
func (t *SecretTrust) Validate() error {
	if t.Mount == "" {
		return fmt.Errorf("trust: mount is required")
	}
	if t.Path == "" {
		return fmt.Errorf("trust: path is required")
	}
	if t.AssignedBy == "" {
		return fmt.Errorf("trust: assigned_by is required")
	}
	if !IsValidTrustLevel(t.Level) {
		return fmt.Errorf("trust: unknown trust level %q", t.Level)
	}
	return nil
}
