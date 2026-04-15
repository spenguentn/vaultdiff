package vault

import (
	"fmt"
	"time"
)

// MaturityLevel represents the maturity stage of a secret.
type MaturityLevel string

const (
	MaturityDraft      MaturityLevel = "draft"
	MaturityCandidate  MaturityLevel = "candidate"
	MaturityStable     MaturityLevel = "stable"
	MaturityDeprecated MaturityLevel = "deprecated"
)

// IsValidMaturityLevel returns true if the given level is recognised.
func IsValidMaturityLevel(l MaturityLevel) bool {
	switch l {
	case MaturityDraft, MaturityCandidate, MaturityStable, MaturityDeprecated:
		return true
	}
	return false
}

// SecretMaturity records the maturity level assigned to a secret.
type SecretMaturity struct {
	Mount      string        `json:"mount"`
	Path       string        `json:"path"`
	Level      MaturityLevel `json:"level"`
	AssignedBy string        `json:"assigned_by"`
	AssignedAt time.Time     `json:"assigned_at"`
}

// FullPath returns the canonical mount+path identifier.
func (m SecretMaturity) FullPath() string {
	return fmt.Sprintf("%s/%s", m.Mount, m.Path)
}

// Validate checks that required fields are present and the level is valid.
func (m SecretMaturity) Validate() error {
	if m.Mount == "" {
		return fmt.Errorf("maturity: mount is required")
	}
	if m.Path == "" {
		return fmt.Errorf("maturity: path is required")
	}
	if m.AssignedBy == "" {
		return fmt.Errorf("maturity: assigned_by is required")
	}
	if !IsValidMaturityLevel(m.Level) {
		return fmt.Errorf("maturity: unknown level %q", m.Level)
	}
	return nil
}
