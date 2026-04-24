package vault

import (
	"fmt"
	"time"
)

// ReputationLevel represents the reputation standing of a secret.
type ReputationLevel string

const (
	ReputationExemplary  ReputationLevel = "exemplary"
	ReputationGood       ReputationLevel = "good"
	ReputationNeutral    ReputationLevel = "neutral"
	ReputationSuspect    ReputationLevel = "suspect"
	ReputationCompromised ReputationLevel = "compromised"
)

// IsValidReputationLevel returns true if the given level is a known reputation level.
func IsValidReputationLevel(level ReputationLevel) bool {
	switch level {
	case ReputationExemplary, ReputationGood, ReputationNeutral, ReputationSuspect, ReputationCompromised:
		return true
	}
	return false
}

// SecretReputation records the reputation standing of a secret at a given path.
type SecretReputation struct {
	Mount      string          `json:"mount"`
	Path       string          `json:"path"`
	Level      ReputationLevel `json:"level"`
	Reason     string          `json:"reason,omitempty"`
	AssessedBy string          `json:"assessed_by"`
	AssessedAt time.Time       `json:"assessed_at"`
}

// FullPath returns the combined mount and path for this reputation record.
func (r SecretReputation) FullPath() string {
	return fmt.Sprintf("%s/%s", r.Mount, r.Path)
}

// Validate checks that the reputation record has all required fields.
func (r SecretReputation) Validate() error {
	if r.Mount == "" {
		return fmt.Errorf("reputation: mount is required")
	}
	if r.Path == "" {
		return fmt.Errorf("reputation: path is required")
	}
	if !IsValidReputationLevel(r.Level) {
		return fmt.Errorf("reputation: unknown level %q", r.Level)
	}
	if r.AssessedBy == "" {
		return fmt.Errorf("reputation: assessed_by is required")
	}
	return nil
}
