package vault

import (
	"fmt"
	"time"
)

// ReadinessStatus represents whether a secret is ready for use.
type ReadinessStatus string

const (
	ReadinessReady    ReadinessStatus = "ready"
	ReadinessNotReady ReadinessStatus = "not_ready"
	ReadinessPending  ReadinessStatus = "pending"
	ReadinessUnknown  ReadinessStatus = "unknown"
)

// IsValidReadinessStatus returns true if the given status is a known readiness value.
func IsValidReadinessStatus(s ReadinessStatus) bool {
	switch s {
	case ReadinessReady, ReadinessNotReady, ReadinessPending, ReadinessUnknown:
		return true
	}
	return false
}

// SecretReadiness records the readiness state of a secret at a point in time.
type SecretReadiness struct {
	Mount      string          `json:"mount"`
	Path       string          `json:"path"`
	Status     ReadinessStatus `json:"status"`
	Reason     string          `json:"reason,omitempty"`
	CheckedAt  time.Time       `json:"checked_at"`
	CheckedBy  string          `json:"checked_by"`
}

// FullPath returns the canonical mount+path identifier.
func (r SecretReadiness) FullPath() string {
	return fmt.Sprintf("%s/%s", r.Mount, r.Path)
}

// IsReady returns true when the status is ReadinessReady.
func (r SecretReadiness) IsReady() bool {
	return r.Status == ReadinessReady
}

// Validate returns an error if the record is missing required fields.
func (r SecretReadiness) Validate() error {
	if r.Mount == "" {
		return fmt.Errorf("readiness: mount is required")
	}
	if r.Path == "" {
		return fmt.Errorf("readiness: path is required")
	}
	if r.CheckedBy == "" {
		return fmt.Errorf("readiness: checked_by is required")
	}
	if !IsValidReadinessStatus(r.Status) {
		return fmt.Errorf("readiness: unknown status %q", r.Status)
	}
	return nil
}
