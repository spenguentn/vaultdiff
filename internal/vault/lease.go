package vault

import (
	"errors"
	"time"
)

// LeaseInfo holds metadata about a Vault secret lease.
type LeaseInfo struct {
	LeaseID   string
	Duration  time.Duration
	Renewable bool
	IssuedAt  time.Time
}

// ExpiresAt returns the absolute expiry time of the lease.
func (l LeaseInfo) ExpiresAt() time.Time {
	return l.IssuedAt.Add(l.Duration)
}

// IsExpired reports whether the lease has passed its expiry time.
func (l LeaseInfo) IsExpired() bool {
	return time.Now().After(l.ExpiresAt())
}

// TTLRemaining returns the duration remaining before the lease expires.
// Returns zero if already expired.
func (l LeaseInfo) TTLRemaining() time.Duration {
	remaining := time.Until(l.ExpiresAt())
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Validate checks that the LeaseInfo contains a usable lease ID and
// a positive duration.
func (l LeaseInfo) Validate() error {
	if l.LeaseID == "" {
		return errors.New("lease: lease ID must not be empty")
	}
	if l.Duration <= 0 {
		return errors.New("lease: duration must be positive")
	}
	return nil
}

// ParseLease constructs a LeaseInfo from raw Vault API response fields.
func ParseLease(leaseID string, leaseDurationSec int, renewable bool) (LeaseInfo, error) {
	if leaseID == "" {
		return LeaseInfo{}, errors.New("lease: empty lease ID")
	}
	if leaseDurationSec <= 0 {
		return LeaseInfo{}, errors.New("lease: duration must be positive")
	}
	return LeaseInfo{
		LeaseID:   leaseID,
		Duration:  time.Duration(leaseDurationSec) * time.Second,
		Renewable: renewable,
		IssuedAt:  time.Now(),
	}, nil
}
