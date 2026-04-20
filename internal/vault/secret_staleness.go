package vault

import (
	"fmt"
	"time"
)

// StalenessLevel represents how stale a secret is considered.
type StalenessLevel string

const (
	StalenessLevelFresh    StalenessLevel = "fresh"
	StalenessLevelWarning  StalenessLevel = "warning"
	StalenessLevelStale    StalenessLevel = "stale"
	StalenessLevelCritical StalenessLevel = "critical"
)

// IsValidStalenessLevel returns true if the level is a known staleness level.
func IsValidStalenessLevel(level StalenessLevel) bool {
	switch level {
	case StalenessLevelFresh, StalenessLevelWarning, StalenessLevelStale, StalenessLevelCritical:
		return true
	}
	return false
}

// SecretStaleness records the staleness evaluation for a secret.
type SecretStaleness struct {
	Mount       string         `json:"mount"`
	Path        string         `json:"path"`
	Level       StalenessLevel `json:"level"`
	LastUpdated time.Time      `json:"last_updated"`
	EvaluatedAt time.Time      `json:"evaluated_at"`
}

// FullPath returns the combined mount and path.
func (s SecretStaleness) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Mount, s.Path)
}

// AgeDays returns the number of days since the secret was last updated.
func (s SecretStaleness) AgeDays() int {
	return int(s.EvaluatedAt.Sub(s.LastUpdated).Hours() / 24)
}

// Validate checks that the staleness record has required fields.
func (s SecretStaleness) Validate() error {
	if s.Mount == "" {
		return fmt.Errorf("staleness: mount is required")
	}
	if s.Path == "" {
		return fmt.Errorf("staleness: path is required")
	}
	if !IsValidStalenessLevel(s.Level) {
		return fmt.Errorf("staleness: unknown level %q", s.Level)
	}
	if s.LastUpdated.IsZero() {
		return fmt.Errorf("staleness: last_updated is required")
	}
	return nil
}

// ComputeStaleness evaluates a secret's staleness given thresholds in days.
// warnDays: days before warning, staleDays: days before stale, critDays: days before critical.
func ComputeStaleness(mount, path string, lastUpdated time.Time, warnDays, staleDays, critDays int) (SecretStaleness, error) {
	if mount == "" {
		return SecretStaleness{}, fmt.Errorf("staleness: mount is required")
	}
	if path == "" {
		return SecretStaleness{}, fmt.Errorf("staleness: path is required")
	}
	now := time.Now().UTC()
	age := int(now.Sub(lastUpdated).Hours() / 24)
	var level StalenessLevel
	switch {
	case age >= critDays:
		level = StalenessLevelCritical
	case age >= staleDays:
		level = StalenessLevelStale
	case age >= warnDays:
		level = StalenessLevelWarning
	default:
		level = StalenessLevelFresh
	}
	return SecretStaleness{
		Mount:       mount,
		Path:        path,
		Level:       level,
		LastUpdated: lastUpdated,
		EvaluatedAt: now,
	}, nil
}
