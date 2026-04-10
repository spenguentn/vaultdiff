// Package compare provides the high-level orchestration layer for vaultdiff.
//
// It coordinates reading secret versions from two Vault instances (or the same
// instance with different paths) and delegates the key-by-key comparison to
// the internal/diff package.
//
// Typical usage:
//
//	engine := compare.NewEngine(leftClient, rightClient)
//	results, err := engine.Run(ctx, leftSource, rightSource)
//
// The Engine is intentionally stateless so that it can be reused across multiple
// Run calls within a single CLI invocation.
package compare
