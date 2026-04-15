package vault

import "fmt"

// ImpactLevel represents the blast radius of a secret change.
type ImpactLevel string

const (
	ImpactLow      ImpactLevel = "low"
	ImpactMedium   ImpactLevel = "medium"
	ImpactHigh     ImpactLevel = "high"
	ImpactCritical ImpactLevel = "critical"
)

// IsValidImpactLevel returns true if the level is a known value.
func IsValidImpactLevel(level ImpactLevel) bool {
	switch level {
	case ImpactLow, ImpactMedium, ImpactHigh, ImpactCritical:
		return true
	}
	return false
}

// SecretImpact records the assessed impact level of a secret at a given path.
type SecretImpact struct {
	Mount       string      `json:"mount"`
	Path        string      `json:"path"`
	Level       ImpactLevel `json:"level"`
	Justification string    `json:"justification,omitempty"`
	AssessedBy  string      `json:"assessed_by"`
}

// FullPath returns the canonical mount+path identifier.
func (s SecretImpact) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Mount, s.Path)
}

// Validate checks that required fields are present and the level is valid.
func (s SecretImpact) Validate() error {
	if s.Mount == "" {
		return fmt.Errorf("impact: mount is required")
	}
	if s.Path == "" {
		return fmt.Errorf("impact: path is required")
	}
	if s.AssessedBy == "" {
		return fmt.Errorf("impact: assessed_by is required")
	}
	if !IsValidImpactLevel(s.Level) {
		return fmt.Errorf("impact: unknown level %q", s.Level)
	}
	return nil
}
