package vault

import (
	"errors"
	"fmt"
	"time"
)

// RetentionUnit represents the time unit for a retention policy.
type RetentionUnit string

const (
	RetentionUnitDays   RetentionUnit = "days"
	RetentionUnitWeeks  RetentionUnit = "weeks"
	RetentionUnitMonths RetentionUnit = "months"
)

// IsValidRetentionUnit returns true if the unit is a known retention unit.
func IsValidRetentionUnit(u RetentionUnit) bool {
	switch u {
	case RetentionUnitDays, RetentionUnitWeeks, RetentionUnitMonths:
		return true
	}
	return false
}

// RetentionPolicy defines how long a secret should be retained before deletion.
type RetentionPolicy struct {
	Mount     string        `json:"mount"`
	Path      string        `json:"path"`
	Duration  int           `json:"duration"`
	Unit      RetentionUnit `json:"unit"`
	CreatedAt time.Time     `json:"created_at"`
	CreatedBy string        `json:"created_by"`
}

// FullPath returns the combined mount and path.
func (r RetentionPolicy) FullPath() string {
	return fmt.Sprintf("%s/%s", r.Mount, r.Path)
}

// ExpiresAt returns the absolute time when the secret should be deleted.
func (r RetentionPolicy) ExpiresAt() time.Time {
	switch r.Unit {
	case RetentionUnitWeeks:
		return r.CreatedAt.AddDate(0, 0, r.Duration*7)
	case RetentionUnitMonths:
		return r.CreatedAt.AddDate(0, r.Duration, 0)
	default:
		return r.CreatedAt.AddDate(0, 0, r.Duration)
	}
}

// IsExpired returns true if the policy's retention window has passed.
func (r RetentionPolicy) IsExpired() bool {
	return time.Now().After(r.ExpiresAt())
}

// Validate checks that the retention policy has all required fields.
func (r RetentionPolicy) Validate() error {
	if r.Mount == "" {
		return errors.New("retention policy: mount is required")
	}
	if r.Path == "" {
		return errors.New("retention policy: path is required")
	}
	if r.Duration <= 0 {
		return errors.New("retention policy: duration must be positive")
	}
	if !IsValidRetentionUnit(r.Unit) {
		return fmt.Errorf("retention policy: unknown unit %q", r.Unit)
	}
	if r.CreatedBy == "" {
		return errors.New("retention policy: created_by is required")
	}
	return nil
}
