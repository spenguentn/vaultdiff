package vault

import (
	"context"
	"fmt"
)

// SecretReader is the subset of client behaviour needed to read a secret.
type SecretReader interface {
	ReadSecret(ctx context.Context, mount, path string) (map[string]interface{}, error)
}

// SecretWriter is the subset of client behaviour needed to write a secret.
type SecretWriter interface {
	WriteSecret(ctx context.Context, mount, path string, data map[string]interface{}) error
}

// SecretCopier copies secrets between paths using a reader and writer.
type SecretCopier struct {
	reader SecretReader
	writer SecretWriter
}

// NewSecretCopier creates a SecretCopier. Panics if either dependency is nil.
func NewSecretCopier(r SecretReader, w SecretWriter) *SecretCopier {
	if r == nil {
		panic("vault: SecretCopier reader must not be nil")
	}
	if w == nil {
		panic("vault: SecretCopier writer must not be nil")
	}
	return &SecretCopier{reader: r, writer: w}
}

// Copy executes a single CopyRequest and returns a CopyResult.
func (c *SecretCopier) Copy(ctx context.Context, req CopyRequest) CopyResult {
	if err := req.Validate(); err != nil {
		return CopyResult{Request: req, Err: err}
	}

	data, err := c.reader.ReadSecret(ctx, req.SourceMount, req.SourcePath)
	if err != nil {
		return CopyResult{Request: req, Err: fmt.Errorf("read: %w", err)}
	}

	if err := c.writer.WriteSecret(ctx, req.DestMount, req.DestPath, data); err != nil {
		return CopyResult{Request: req, Err: fmt.Errorf("write: %w", err)}
	}

	return CopyResult{Request: req}
}

// CopyPlanAll executes every request in the plan and returns all results.
func (c *SecretCopier) CopyPlanAll(ctx context.Context, plan *CopyPlan) []CopyResult {
	results := make([]CopyResult, 0, plan.Len())
	for _, req := range plan.Requests() {
		results = append(results, c.Copy(ctx, req))
	}
	return results
}
