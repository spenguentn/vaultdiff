package vault

import (
	"fmt"
	"time"
)

// PriorityLevel represents the urgency level of a secret.
type PriorityLevel string

const (
	PriorityLow      PriorityLevel = "low"
	PriorityMedium   PriorityLevel = "medium"
	PriorityHigh     PriorityLevel = "high"
	PriorityCritical PriorityLevel = "critical"
)

// IsValidPriorityLevel reports whether the given level is recognized.
func IsValidPriorityLevel(level PriorityLevel) bool {
	switch level {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	}
	return false
}

// SecretPriority associates a priority level with a secret path.
type SecretPriority struct {
	Mount      string        `json:"mount"`
	Path       string        `json:"path"`
	Level      PriorityLevel `json:"level"`
	Reason     string        `json:"reason,omitempty"`
	AssignedBy string        `json:"assigned_by"`
	AssignedAt time.Time     `json:"assigned_at"`
}

// FullPath returns the canonical mount+path key.
func (p SecretPriority) FullPath() string {
	return p.Mount + "/" + p.Path
}

// Validate checks that the priority record is well-formed.
func (p SecretPriority) Validate() error {
	if p.Mount == "" {
		return fmt.Errorf("secret priority: mount is required")
	}
	if p.Path == "" {
		return fmt.Errorf("secret priority: path is required")
	}
	if !IsValidPriorityLevel(p.Level) {
		return fmt.Errorf("secret priority: unknown level %q", p.Level)
	}
	if p.AssignedBy == "" {
		return fmt.Errorf("secret priority: assigned_by is required")
	}
	return nil
}
