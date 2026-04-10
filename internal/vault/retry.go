package vault

import (
	"errors"
	"time"
)

// RetryConfig controls retry behaviour for Vault API calls.
type RetryConfig struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the first retry.
	InitialDelay time.Duration
	// MaxDelay caps the exponential back-off.
	MaxDelay time.Duration
}

// DefaultRetryConfig returns a RetryConfig suitable for most Vault operations.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 200 * time.Millisecond,
		MaxDelay:     2 * time.Second,
	}
}

// Validate returns an error if the RetryConfig contains invalid values.
func (r RetryConfig) Validate() error {
	if r.MaxAttempts < 1 {
		return errors.New("retry: MaxAttempts must be at least 1")
	}
	if r.InitialDelay < 0 {
		return errors.New("retry: InitialDelay must be non-negative")
	}
	if r.MaxDelay < r.InitialDelay {
		return errors.New("retry: MaxDelay must be >= InitialDelay")
	}
	return nil
}

// Do executes fn up to cfg.MaxAttempts times, backing off exponentially between
// attempts. It returns the first nil error or the last non-nil error.
func Do(cfg RetryConfig, fn func() error) error {
	delay := cfg.InitialDelay
	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if attempt < cfg.MaxAttempts-1 {
			time.Sleep(delay)
			delay *= 2
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}
	}
	return lastErr
}
