package vault

import (
	"errors"
	"fmt"
	"time"
)

// DeprecationStatus represents the deprecation state of a secret.
type DeprecationStatus string

const (
	DeprecationStatusActive     DeprecationStatus = "active"
	DeprecationStatusDeprecated DeprecationStatus = "deprecated"
	DeprecationStatusSunset     DeprecationStatus = "sunset"
)

// IsValidDeprecationStatus reports whether s is a known deprecation status.
func IsValidDeprecationStatus(s DeprecationStatus) bool {
	switch s {
	case DeprecationStatusActive, DeprecationStatusDeprecated, DeprecationStatusSunset:
		return true
	}
	return false
}

// SecretDeprecation records the deprecation notice for a secret.
type SecretDeprecation struct {
	Mount         string            `json:"mount"`
	Path          string            `json:"path"`
	Status        DeprecationStatus `json:"status"`
	Reason        string            `json:"reason"`
	DeprecatedBy  string            `json:"deprecated_by"`
	SunsetAt      *time.Time        `json:"sunset_at,omitempty"`
	DeprecatedAt  time.Time         `json:"deprecated_at"`
	ReplacementPath string          `json:"replacement_path,omitempty"`
}

// FullPath returns the canonical mount+path identifier.
func (d SecretDeprecation) FullPath() string {
	return fmt.Sprintf("%s/%s", d.Mount, d.Path)
}

// IsSunset reports whether the deprecation has reached sunset.
func (d SecretDeprecation) IsSunset() bool {
	if d.SunsetAt == nil {
		return false
	}
	return time.Now().After(*d.SunsetAt)
}

// Validate returns an error if the deprecation record is incomplete or invalid.
func (d SecretDeprecation) Validate() error {
	if d.Mount == "" {
		return errors.New("deprecation: mount is required")
	}
	if d.Path == "" {
		return errors.New("deprecation: path is required")
	}
	if !IsValidDeprecationStatus(d.Status) {
		return fmt.Errorf("deprecation: unknown status %q", d.Status)
	}
	if d.DeprecatedBy == "" {
		return errors.New("deprecation: deprecated_by is required")
	}
	if d.Reason == "" {
		return errors.New("deprecation: reason is required")
	}
	return nil
}
