package vault

import (
	"errors"
	"fmt"
	"time"
)

// ReminderFrequency defines how often a reminder should trigger.
type ReminderFrequency string

const (
	ReminderFrequencyOnce    ReminderFrequency = "once"
	ReminderFrequencyWeekly  ReminderFrequency = "weekly"
	ReminderFrequencyMonthly ReminderFrequency = "monthly"
)

// SecretReminder represents a scheduled reminder attached to a secret.
type SecretReminder struct {
	Mount     string            `json:"mount"`
	Path      string            `json:"path"`
	Message   string            `json:"message"`
	RemindAt  time.Time         `json:"remind_at"`
	Frequency ReminderFrequency `json:"frequency"`
	CreatedBy string            `json:"created_by"`
	CreatedAt time.Time         `json:"created_at"`
}

// FullPath returns the canonical mount+path identifier.
func (r SecretReminder) FullPath() string {
	return fmt.Sprintf("%s/%s", r.Mount, r.Path)
}

// IsDue reports whether the reminder is due relative to the given time.
func (r SecretReminder) IsDue(now time.Time) bool {
	return !now.Before(r.RemindAt)
}

// Validate ensures the reminder has all required fields.
func (r SecretReminder) Validate() error {
	if r.Mount == "" {
		return errors.New("reminder: mount is required")
	}
	if r.Path == "" {
		return errors.New("reminder: path is required")
	}
	if r.Message == "" {
		return errors.New("reminder: message is required")
	}
	if r.RemindAt.IsZero() {
		return errors.New("reminder: remind_at is required")
	}
	if r.CreatedBy == "" {
		return errors.New("reminder: created_by is required")
	}
	return nil
}

// IsValidFrequency reports whether f is a known ReminderFrequency.
func IsValidFrequency(f ReminderFrequency) bool {
	switch f {
	case ReminderFrequencyOnce, ReminderFrequencyWeekly, ReminderFrequencyMonthly:
		return true
	}
	return false
}
