package vault

import (
	"fmt"
	"strings"
	"time"
)

// RiskLevel represents the assessed risk level of a secret.
type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

// IsValidRiskLevel returns true if the given level is a known risk level.
func IsValidRiskLevel(level RiskLevel) bool {
	switch level {
	case RiskLevelLow, RiskLevelMedium, RiskLevelHigh, RiskLevelCritical:
		return true
	}
	return false
}

// SecretRisk holds a risk assessment record for a secret.
type SecretRisk struct {
	Mount      string    `json:"mount"`
	Path       string    `json:"path"`
	Level      RiskLevel `json:"level"`
	Reason     string    `json:"reason"`
	AssessedBy string    `json:"assessed_by"`
	AssessedAt time.Time `json:"assessed_at"`
}

// FullPath returns the canonical mount+path identifier.
func (r *SecretRisk) FullPath() string {
	return fmt.Sprintf("%s/%s", strings.Trim(r.Mount, "/"), strings.Trim(r.Path, "/"))
}

// Validate returns an error if the risk record is incomplete or invalid.
func (r *SecretRisk) Validate() error {
	if r.Mount == "" {
		return fmt.Errorf("risk: mount is required")
	}
	if r.Path == "" {
		return fmt.Errorf("risk: path is required")
	}
	if !IsValidRiskLevel(r.Level) {
		return fmt.Errorf("risk: unknown level %q", r.Level)
	}
	if r.AssessedBy == "" {
		return fmt.Errorf("risk: assessed_by is required")
	}
	return nil
}
