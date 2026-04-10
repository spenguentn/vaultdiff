package vault

import (
	"context"
	"fmt"
	"time"
)

// RollbackRequest describes a request to restore a secret to a prior version.
type RollbackRequest struct {
	Mount   string
	Path    string
	Version int
}

// RollbackResult holds the outcome of a rollback operation.
type RollbackResult struct {
	Mount      string
	Path       string
	Version    int
	RolledBack bool
	RolledAt   time.Time
	Err        error
}

// IsSuccess returns true when the rollback completed without error.
func (r RollbackResult) IsSuccess() bool {
	return r.RolledBack && r.Err == nil
}

// String returns a human-readable summary of the result.
func (r RollbackResult) String() string {
	if r.Err != nil {
		return fmt.Sprintf("rollback %s/%s@v%d FAILED: %v", r.Mount, r.Path, r.Version, r.Err)
	}
	return fmt.Sprintf("rollback %s/%s@v%d OK at %s", r.Mount, r.Path, r.Version, r.RolledAt.Format(time.RFC3339))
}

// Rollbacker can restore a KV secret to a previous version.
type Rollbacker struct {
	client *Client
}

// NewRollbacker creates a Rollbacker using the provided Vault client.
func NewRollbacker(c *Client) *Rollbacker {
	if c == nil {
		panic("vault: NewRollbacker requires a non-nil client")
	}
	return &Rollbacker{client: c}
}

// Rollback restores the secret at req.Path to req.Version using the KV v2
// undelete + restore pattern (write to <mount>/data/<path> with the version's
// data is not exposed by the API, so we use the official "rollback" approach:
// read the target version then write it as a new version).
func (rb *Rollbacker) Rollback(ctx context.Context, req RollbackRequest) RollbackResult {
	result := RollbackResult{Mount: req.Mount, Path: req.Path, Version: req.Version}

	if req.Version < 1 {
		result.Err = fmt.Errorf("invalid version %d: must be >= 1", req.Version)
		return result
	}

	logical := rb.client.Logical()

	// Read the target version.
	dataPath := fmt.Sprintf("%s/data/%s", req.Mount, req.Path)
	secret, err := logical.ReadWithDataWithContext(ctx, dataPath, map[string][]string{
		"version": {fmt.Sprintf("%d", req.Version)},
	})
	if err != nil {
		result.Err = fmt.Errorf("read version %d: %w", req.Version, err)
		return result
	}
	if secret == nil || secret.Data == nil {
		result.Err = fmt.Errorf("version %d not found at %s", req.Version, req.Path)
		return result
	}

	kvData, _ := secret.Data["data"].(map[string]interface{})
	_, err = logical.WriteWithContext(ctx, dataPath, map[string]interface{}{"data": kvData})
	if err != nil {
		result.Err = fmt.Errorf("write rollback data: %w", err)
		return result
	}

	result.RolledBack = true
	result.RolledAt = time.Now().UTC()
	return result
}
