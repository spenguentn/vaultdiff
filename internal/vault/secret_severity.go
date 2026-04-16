package vault

import "fmt"

// SeverityLevel represents the severity of a secret issue.
type SeverityLevel string

const (
	SeverityLow      SeverityLevel = "low"
	SeverityMedium   SeverityLevel = "medium"
	SeverityHigh     SeverityLevel = "high"
	SeverityCritical SeverityLevel = "critical"
)

var validSeverityLevels = map[SeverityLevel]bool{
	SeverityLow:      true,
	SeverityMedium:   true,
	SeverityHigh:     true,
	SeverityCritical: true,
}

// IsValidSeverityLevel returns true if the given level is recognised.
func IsValidSeverityLevel(level SeverityLevel) bool {
	return validSeverityLevels[level]
}

// SecretSeverity associates a severity level with a secret path.
type SecretSeverity struct {
	Mount    string        `json:"mount"`
	Path     string        `json:"path"`
	Level    SeverityLevel `json:"level"`
	Reason   string        `json:"reason,omitempty"`
	AssignedBy string      `json:"assigned_by"`
}

// FullPath returns the canonical mount+path string.
func (s SecretSeverity) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Mount, s.Path)
}

// Validate checks that the SecretSeverity is well-formed.
func (s SecretSeverity) Validate() error {
	if s.Mount == "" {
		return fmt.Errorf("severity: mount is required")
	}
	if s.Path == "" {
		return fmt.Errorf("severity: path is required")
	}
	if !IsValidSeverityLevel(s.Level) {
		return fmt.Errorf("severity: unknown level %q", s.Level)
	}
	if s.AssignedBy == "" {
		return fmt.Errorf("severity: assigned_by is required")
	}
	return nil
}
