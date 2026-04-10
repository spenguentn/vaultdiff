package vault

import (
	"context"
	"fmt"
	"time"
)

// HealthStatus represents the result of a Vault health check.
type HealthStatus struct {
	Address     string
	Initialized bool
	Sealed      bool
	Standby     bool
	Version     string
	CheckedAt   time.Time
}

// IsHealthy returns true when Vault is initialized, unsealed, and active.
func (h HealthStatus) IsHealthy() bool {
	return h.Initialized && !h.Sealed && !h.Standby
}

// String returns a human-readable summary of the health status.
func (h HealthStatus) String() string {
	if h.IsHealthy() {
		return fmt.Sprintf("%s — healthy (v%s)", h.Address, h.Version)
	}
	if h.Sealed {
		return fmt.Sprintf("%s — sealed", h.Address)
	}
	if !h.Initialized {
		return fmt.Sprintf("%s — not initialized", h.Address)
	}
	return fmt.Sprintf("%s — standby", h.Address)
}

// CheckHealth queries the Vault sys/health endpoint and returns a HealthStatus.
// It uses the provided context for cancellation and timeout control.
func CheckHealth(ctx context.Context, c *Client) (HealthStatus, error) {
	if c == nil {
		return HealthStatus{}, fmt.Errorf("vault: client must not be nil")
	}

	resp, err := c.inner.Sys().HealthWithContext(ctx)
	if err != nil {
		return HealthStatus{}, fmt.Errorf("vault: health check failed for %s: %w", c.Address(), err)
	}

	return HealthStatus{
		Address:     c.Address(),
		Initialized: resp.Initialized,
		Sealed:      resp.Sealed,
		Standby:     resp.Standby,
		Version:     resp.Version,
		CheckedAt:   time.Now().UTC(),
	}, nil
}
