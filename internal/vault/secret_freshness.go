package vault

import (
	"fmt"
	"time"
)

// FreshnessStatus represents how up-to-date a secret is relative to its expected update cadence.
type FreshnessStatus string

const (
	FreshnessFresh   FreshnessStatus = "fresh"
	FreshnessStale   FreshnessStatus = "stale"
	FreshnessExpired FreshnessStatus = "expired"
	FreshnessUnknown FreshnessStatus = "unknown"
)

// IsValidFreshnessStatus returns true if the given status is a known freshness value.
func IsValidFreshnessStatus(s FreshnessStatus) bool {
	switch s {
	case FreshnessFresh, FreshnessStale, FreshnessExpired, FreshnessUnknown:
		return true
	}
	return false
}

// SecretFreshness tracks how recently a secret was updated relative to an expected cadence.
type SecretFreshness struct {
	Mount       string          `json:"mount"`
	Path        string          `json:"path"`
	LastUpdated time.Time       `json:"last_updated"`
	MaxAge      time.Duration   `json:"max_age"`
	Status      FreshnessStatus `json:"status"`
	CheckedAt   time.Time       `json:"checked_at"`
}

// FullPath returns the combined mount and path for this freshness record.
func (f *SecretFreshness) FullPath() string {
	return fmt.Sprintf("%s/%s", f.Mount, f.Path)
}

// Validate returns an error if the freshness record is missing required fields.
func (f *SecretFreshness) Validate() error {
	if f.Mount == "" {
		return fmt.Errorf("freshness: mount is required")
	}
	if f.Path == "" {
		return fmt.Errorf("freshness: path is required")
	}
	if f.MaxAge <= 0 {
		return fmt.Errorf("freshness: max_age must be positive")
	}
	if f.LastUpdated.IsZero() {
		return fmt.Errorf("freshness: last_updated must be set")
	}
	if !IsValidFreshnessStatus(f.Status) {
		return fmt.Errorf("freshness: unknown status %q", f.Status)
	}
	return nil
}

// ComputeFreshness evaluates a secret's freshness based on its last update time and max age.
func ComputeFreshness(mount, path string, lastUpdated time.Time, maxAge time.Duration) (*SecretFreshness, error) {
	if mount == "" {
		return nil, fmt.Errorf("freshness: mount is required")
	}
	if path == "" {
		return nil, fmt.Errorf("freshness: path is required")
	}
	if maxAge <= 0 {
		return nil, fmt.Errorf("freshness: max_age must be positive")
	}

	now := time.Now().UTC()
	age := now.Sub(lastUpdated)

	var status FreshnessStatus
	switch {
	case lastUpdated.IsZero():
		status = FreshnessUnknown
	case age <= maxAge/2:
		status = FreshnessFresh
	case age <= maxAge:
		status = FreshnessStale
	default:
		status = FreshnessExpired
	}

	return &SecretFreshness{
		Mount:       mount,
		Path:        path,
		LastUpdated: lastUpdated,
		MaxAge:      maxAge,
		Status:      status,
		CheckedAt:   now,
	}, nil
}
