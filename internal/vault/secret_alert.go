package vault

import (
	"errors"
	"fmt"
	"time"
)

// AlertSeverity represents the urgency level of a secret alert.
type AlertSeverity string

const (
	AlertSeverityLow      AlertSeverity = "low"
	AlertSeverityMedium   AlertSeverity = "medium"
	AlertSeverityHigh     AlertSeverity = "high"
	AlertSeverityCritical AlertSeverity = "critical"
)

// IsValidSeverity returns true if s is a recognised AlertSeverity value.
func IsValidSeverity(s AlertSeverity) bool {
	switch s {
	case AlertSeverityLow, AlertSeverityMedium, AlertSeverityHigh, AlertSeverityCritical:
		return true
	}
	return false
}

// SecretAlert represents a triggered alert for a secret at a given path.
type SecretAlert struct {
	Mount     string        `json:"mount"`
	Path      string        `json:"path"`
	Message   string        `json:"message"`
	Severity  AlertSeverity `json:"severity"`
	Triggered time.Time     `json:"triggered"`
	Actor     string        `json:"actor"`
}

// FullPath returns the combined mount and path for the alerted secret.
func (a SecretAlert) FullPath() string {
	return fmt.Sprintf("%s/%s", a.Mount, a.Path)
}

// Validate checks that the alert contains the required fields.
func (a SecretAlert) Validate() error {
	if a.Mount == "" {
		return errors.New("alert: mount is required")
	}
	if a.Path == "" {
		return errors.New("alert: path is required")
	}
	if a.Message == "" {
		return errors.New("alert: message is required")
	}
	if !IsValidSeverity(a.Severity) {
		return fmt.Errorf("alert: invalid severity %q", a.Severity)
	}
	if a.Actor == "" {
		return errors.New("alert: actor is required")
	}
	return nil
}
