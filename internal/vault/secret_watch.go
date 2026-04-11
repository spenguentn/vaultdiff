package vault

import (
	"errors"
	"time"
)

// SecretWatchEvent represents a detected change on a watched secret path.
type SecretWatchEvent struct {
	Mount     string
	Path      string
	Version   int
	ChangedAt time.Time
	PrevHash  string
	CurrHash  string
}

// IsChanged returns true when the content hash has changed.
func (e SecretWatchEvent) IsChanged() bool {
	return e.PrevHash != e.CurrHash
}

// SecretWatchConfig holds configuration for watching a single secret.
type SecretWatchConfig struct {
	Mount    string
	Path     string
	Interval time.Duration
	OnChange func(SecretWatchEvent)
}

// Validate checks that the watch config is complete and usable.
func (c SecretWatchConfig) Validate() error {
	if c.Mount == "" {
		return errors.New("secret watch: mount is required")
	}
	if c.Path == "" {
		return errors.New("secret watch: path is required")
	}
	if c.Interval <= 0 {
		return errors.New("secret watch: interval must be positive")
	}
	if c.OnChange == nil {
		return errors.New("secret watch: OnChange handler is required")
	}
	return nil
}

// DefaultWatchInterval is the default polling interval for secret watches.
const DefaultWatchInterval = 30 * time.Second
