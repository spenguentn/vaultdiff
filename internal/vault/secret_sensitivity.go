package vault

import "fmt"

// SensitivityLevel represents how sensitive a secret is.
type SensitivityLevel string

const (
	SensitivityLow      SensitivityLevel = "low"
	SensitivityMedium   SensitivityLevel = "medium"
	SensitivityHigh     SensitivityLevel = "high"
	SensitivityCritical SensitivityLevel = "critical"
)

var validSensitivityLevels = map[SensitivityLevel]bool{
	SensitivityLow:      true,
	SensitivityMedium:   true,
	SensitivityHigh:     true,
	SensitivityCritical: true,
}

// IsValidSensitivityLevel returns true if the level is a known value.
func IsValidSensitivityLevel(level SensitivityLevel) bool {
	return validSensitivityLevels[level]
}

// SecretSensitivity records the sensitivity classification for a secret.
type SecretSensitivity struct {
	Mount       string           `json:"mount"`
	Path        string           `json:"path"`
	Level       SensitivityLevel `json:"level"`
	SetBy       string           `json:"set_by"`
	Justification string         `json:"justification,omitempty"`
}

// FullPath returns the combined mount and path.
func (s SecretSensitivity) FullPath() string {
	return s.Mount + "/" + s.Path
}

// Validate checks that required fields are present and the level is valid.
func (s SecretSensitivity) Validate() error {
	if s.Mount == "" {
		return fmt.Errorf("sensitivity: mount is required")
	}
	if s.Path == "" {
		return fmt.Errorf("sensitivity: path is required")
	}
	if s.SetBy == "" {
		return fmt.Errorf("sensitivity: set_by is required")
	}
	if !IsValidSensitivityLevel(s.Level) {
		return fmt.Errorf("sensitivity: unknown level %q", s.Level)
	}
	return nil
}

// RequiresRedaction returns true when the level warrants automatic masking.
func (s SecretSensitivity) RequiresRedaction() bool {
	return s.Level == SensitivityHigh || s.Level == SensitivityCritical
}
