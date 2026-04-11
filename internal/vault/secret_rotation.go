package vault

import (
	"errors"
	"fmt"
	"time"
)

// RotationPolicy defines how often a secret should be rotated.
type RotationPolicy struct {
	Mount    string
	Path     string
	Interval time.Duration
	LastRotatedAt time.Time
	Enabled  bool
}

// Validate checks that the rotation policy is well-formed.
func (r RotationPolicy) Validate() error {
	if r.Mount == "" {
		return errors.New("rotation policy: mount is required")
	}
	if r.Path == "" {
		return errors.New("rotation policy: path is required")
	}
	if r.Interval <= 0 {
		return errors.New("rotation policy: interval must be positive")
	}
	return nil
}

// IsDue reports whether the secret is due for rotation based on the policy.
func (r RotationPolicy) IsDue() bool {
	if !r.Enabled {
		return false
	}
	if r.LastRotatedAt.IsZero() {
		return true
	}
	return time.Since(r.LastRotatedAt) >= r.Interval
}

// NextRotationAt returns the time at which the next rotation is due.
func (r RotationPolicy) NextRotationAt() time.Time {
	if r.LastRotatedAt.IsZero() {
		return time.Time{}
	}
	return r.LastRotatedAt.Add(r.Interval)
}

// RotationResult holds the outcome of a rotation operation.
type RotationResult struct {
	Mount     string
	Path      string
	RotatedAt time.Time
	Err       error
}

// IsSuccess reports whether the rotation completed without error.
func (r RotationResult) IsSuccess() bool {
	return r.Err == nil
}

// String returns a human-readable summary of the rotation result.
func (r RotationResult) String() string {
	if r.Err != nil {
		return fmt.Sprintf("rotation failed for %s/%s: %v", r.Mount, r.Path, r.Err)
	}
	return fmt.Sprintf("rotated %s/%s at %s", r.Mount, r.Path, r.RotatedAt.Format(time.RFC3339))
}
