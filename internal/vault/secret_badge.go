package vault

import (
	"fmt"
	"time"
)

// BadgeType represents the category of a secret badge.
type BadgeType string

const (
	BadgeTypeCompliant  BadgeType = "compliant"
	BadgeTypeRotated    BadgeType = "rotated"
	BadgeTypeVerified   BadgeType = "verified"
	BadgeTypeDeprecated BadgeType = "deprecated"
)

// IsValidBadgeType returns true if the given badge type is recognized.
func IsValidBadgeType(t BadgeType) bool {
	switch t {
	case BadgeTypeCompliant, BadgeTypeRotated, BadgeTypeVerified, BadgeTypeDeprecated:
		return true
	}
	return false
}

// SecretBadge represents a badge awarded to a secret path for a specific reason.
type SecretBadge struct {
	Mount     string    `json:"mount"`
	Path      string    `json:"path"`
	Badge     BadgeType `json:"badge"`
	AwardedBy string    `json:"awarded_by"`
	AwardedAt time.Time `json:"awarded_at"`
	Note      string    `json:"note,omitempty"`
}

// FullPath returns the fully qualified secret path including mount.
func (b SecretBadge) FullPath() string {
	return fmt.Sprintf("%s/%s", b.Mount, b.Path)
}

// Validate returns an error if the badge is missing required fields or has an invalid type.
func (b SecretBadge) Validate() error {
	if b.Mount == "" {
		return fmt.Errorf("badge: mount is required")
	}
	if b.Path == "" {
		return fmt.Errorf("badge: path is required")
	}
	if !IsValidBadgeType(b.Badge) {
		return fmt.Errorf("badge: unknown badge type %q", b.Badge)
	}
	if b.AwardedBy == "" {
		return fmt.Errorf("badge: awarded_by is required")
	}
	if b.AwardedAt.IsZero() {
		return fmt.Errorf("badge: awarded_at is required")
	}
	return nil
}
