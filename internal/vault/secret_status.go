package vault

import (
	"fmt"
	"time"
)

// SecretStatusStage represents the current operational status of a secret.
type SecretStatusStage string

const (
	StatusActive    SecretStatusStage = "active"
	StatusDeprecated SecretStatusStage = "deprecated"
	StatusPendingRotation SecretStatusStage = "pending_rotation"
	StatusExpired   SecretStatusStage = "expired"
	StatusRevoked   SecretStatusStage = "revoked"
)

// IsValidStatus returns true if the given stage is a known SecretStatusStage.
func IsValidStatus(s SecretStatusStage) bool {
	switch s {
	case StatusActive, StatusDeprecated, StatusPendingRotation, StatusExpired, StatusRevoked:
		return true
	}
	return false
}

// SecretStatus records the current status of a secret at a given path.
type SecretStatus struct {
	Mount     string            `json:"mount"`
	Path      string            `json:"path"`
	Stage     SecretStatusStage `json:"stage"`
	Reason    string            `json:"reason,omitempty"`
	UpdatedBy string            `json:"updated_by"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// FullPath returns the canonical mount+path identifier.
func (s SecretStatus) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Mount, s.Path)
}

// Validate returns an error if the SecretStatus is missing required fields.
func (s SecretStatus) Validate() error {
	if s.Mount == "" {
		return fmt.Errorf("secret status: mount is required")
	}
	if s.Path == "" {
		return fmt.Errorf("secret status: path is required")
	}
	if s.UpdatedBy == "" {
		return fmt.Errorf("secret status: updated_by is required")
	}
	if !IsValidStatus(s.Stage) {
		return fmt.Errorf("secret status: unknown stage %q", s.Stage)
	}
	return nil
}
