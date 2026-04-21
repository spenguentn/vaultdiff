package vault

import (
	"errors"
	"fmt"
	"time"
)

// IsValidRevocationReason returns true if the given reason is a known revocation reason.
func IsValidRevocationReason(reason string) bool {
	switch reason {
	case "compromised", "expired", "rotated", "decommissioned", "policy_violation":
		return true
	}
	return false
}

// RevocationRecord captures the details of a secret revocation event.
type RevocationRecord struct {
	Mount		string    `json:"mount"`
	Path		string    `json:"path"`
	RevokedBy	string    `json:"revoked_by"`
	Reason		string    `json:"reason"`
	Note		string    `json:"note,omitempty"`
	RevokedAt	time.Time `json:"revoked_at"`
	ExpiresAt	*time.Time `json:"expires_at,omitempty"`
}

// FullPath returns the canonical vault path for this revocation record.
func (r RevocationRecord) FullPath() string {
	return fmt.Sprintf("%s/%s", r.Mount, r.Path)
}

// IsExpired returns true if the revocation has a defined expiry that has passed.
func (r RevocationRecord) IsExpired() bool {
	if r.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*r.ExpiresAt)
}

// Validate returns an error if the revocation record is incomplete or invalid.
func (r RevocationRecord) Validate() error {
	if r.Mount == "" {
		return errors.New("revocation: mount is required")
	}
	if r.Path == "" {
		return errors.New("revocation: path is required")
	}
	if r.RevokedBy == "" {
		return errors.New("revocation: revokedBy is required")
	}
	if !IsValidRevocationReason(r.Reason) {
		return fmt.Errorf("revocation: unknown reason %q", r.Reason)
	}
	if r.RevokedAt.IsZero() {
		return errors.New("revocation: revokedAt must be set")
	}
	return nil
}
