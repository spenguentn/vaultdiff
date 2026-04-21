package vault

import (
	"fmt"
	"time"
)

// CoverageStatus represents the documentation/metadata coverage level of a secret.
type CoverageStatus string

const (
	CoverageStatusFull    CoverageStatus = "full"
	CoverageStatusPartial CoverageStatus = "partial"
	CoverageStatusNone    CoverageStatus = "none"
)

// IsValidCoverageStatus returns true if the given status is a known coverage status.
func IsValidCoverageStatus(s CoverageStatus) bool {
	switch s {
	case CoverageStatusFull, CoverageStatusPartial, CoverageStatusNone:
		return true
	}
	return false
}

// SecretCoverage records the metadata/documentation coverage state of a secret.
type SecretCoverage struct {
	Mount      string         `json:"mount"`
	Path       string         `json:"path"`
	Status     CoverageStatus `json:"status"`
	Score      int            `json:"score"` // 0–100
	AssessedBy string         `json:"assessed_by"`
	AssessedAt time.Time      `json:"assessed_at"`
	Notes      string         `json:"notes,omitempty"`
}

// FullPath returns the canonical mount+path identifier.
func (c SecretCoverage) FullPath() string {
	return fmt.Sprintf("%s/%s", c.Mount, c.Path)
}

// Validate returns an error if the coverage record is incomplete or invalid.
func (c SecretCoverage) Validate() error {
	if c.Mount == "" {
		return fmt.Errorf("coverage: mount is required")
	}
	if c.Path == "" {
		return fmt.Errorf("coverage: path is required")
	}
	if !IsValidCoverageStatus(c.Status) {
		return fmt.Errorf("coverage: invalid status %q", c.Status)
	}
	if c.Score < 0 || c.Score > 100 {
		return fmt.Errorf("coverage: score must be between 0 and 100, got %d", c.Score)
	}
	if c.AssessedBy == "" {
		return fmt.Errorf("coverage: assessed_by is required")
	}
	if c.AssessedAt.IsZero() {
		return fmt.Errorf("coverage: assessed_at is required")
	}
	return nil
}
