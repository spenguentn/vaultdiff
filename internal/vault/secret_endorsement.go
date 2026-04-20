package vault

import (
	"errors"
	"fmt"
	"time"
)

// IsValidEndorsementStatus returns true if the given status is a known endorsement status.
func IsValidEndorsementStatus(s string) bool {
	switch s {
	case "pending", "approved", "rejected", "revoked":
		return true
	}
	return false
}

// SecretEndorsement represents a peer endorsement record for a secret.
type SecretEndorsement struct {
	Mount       string    `json:"mount"`
	Path        string    `json:"path"`
	EndorsedBy  string    `json:"endorsed_by"`
	Status      string    `json:"status"`
	Comment     string    `json:"comment,omitempty"`
	EndorsedAt  time.Time `json:"endorsed_at"`
}

// FullPath returns the canonical mount+path identifier.
func (e *SecretEndorsement) FullPath() string {
	return fmt.Sprintf("%s/%s", e.Mount, e.Path)
}

// Validate checks that the endorsement record is complete and valid.
func (e *SecretEndorsement) Validate() error {
	if e.Mount == "" {
		return errors.New("endorsement: mount is required")
	}
	if e.Path == "" {
		return errors.New("endorsement: path is required")
	}
	if e.EndorsedBy == "" {
		return errors.New("endorsement: endorsed_by is required")
	}
	if !IsValidEndorsementStatus(e.Status) {
		return fmt.Errorf("endorsement: unknown status %q", e.Status)
	}
	return nil
}

// IsApproved returns true when the endorsement has been approved.
func (e *SecretEndorsement) IsApproved() bool {
	return e.Status == "approved"
}
