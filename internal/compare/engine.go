// Package compare orchestrates fetching secrets from two Vault paths and
// producing a diff result suitable for display and audit logging.
package compare

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// Source describes one side of a comparison (environment + path).
type Source struct {
	Environment string
	Mount       string
	SecretPath  string
	Version     int // 0 means latest
}

// Engine fetches and compares secrets from two Vault clients.
type Engine struct {
	Left  *vault.Client
	Right *vault.Client
}

// NewEngine creates an Engine with two pre-configured Vault clients.
func NewEngine(left, right *vault.Client) *Engine {
	return &Engine{Left: left, Right: right}
}

// Run fetches both secret versions and returns the diff results.
func (e *Engine) Run(ctx context.Context, left, right Source) ([]diff.Result, error) {
	leftSecret, err := e.fetchSecret(ctx, e.Left, left)
	if err != nil {
		return nil, fmt.Errorf("fetching left secret (%s): %w", left.Environment, err)
	}

	rightSecret, err := e.fetchSecret(ctx, e.Right, right)
	if err != nil {
		return nil, fmt.Errorf("fetching right secret (%s): %w", right.Environment, err)
	}

	return diff.Compare(leftSecret.Data, rightSecret.Data), nil
}

func (e *Engine) fetchSecret(ctx context.Context, c *vault.Client, src Source) (*vault.SecretVersion, error) {
	if src.Version > 0 {
		return c.ReadSecretVersion(ctx, src.Mount, src.SecretPath, src.Version)
	}
	return c.ReadSecret(ctx, src.Mount, src.SecretPath)
}
