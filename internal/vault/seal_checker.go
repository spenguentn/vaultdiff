package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// SealChecker queries a Vault client for the current seal status.
type SealChecker struct {
	client  *vaultapi.Client
	timeout time.Duration
}

// NewSealChecker creates a SealChecker for the given Vault client.
// If timeout is zero, a default of 5 seconds is used.
func NewSealChecker(client *vaultapi.Client, timeout time.Duration) (*SealChecker, error) {
	if client == nil {
		return nil, fmt.Errorf("seal checker: client must not be nil")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &SealChecker{client: client, timeout: timeout}, nil
}

// Check fetches the current seal status from Vault.
func (sc *SealChecker) Check(ctx context.Context) (SealInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, sc.timeout)
	defer cancel()

	resp, err := sc.client.RawRequestWithContext(ctx,
		sc.client.NewRequest("GET", "/v1/sys/seal-status"),
	)
	if err != nil {
		return SealInfo{}, fmt.Errorf("seal checker: request failed: %w", err)
	}
	defer resp.Body.Close()

	var raw map[string]any
	if err := resp.DecodeJSON(&raw); err != nil {
		return SealInfo{}, fmt.Errorf("seal checker: decode failed: %w", err)
	}
	return ParseSealInfo(raw)
}

// MustBeUnsealed returns an error if the Vault instance is sealed or uninitialized.
func (sc *SealChecker) MustBeUnsealed(ctx context.Context) error {
	info, err := sc.Check(ctx)
	if err != nil {
		return err
	}
	if !info.Initialized {
		return fmt.Errorf("seal checker: vault is not initialized")
	}
	if info.Sealed {
		return fmt.Errorf("seal checker: vault is sealed (progress %d/%d)", info.Progress, info.Threshold)
	}
	return nil
}
