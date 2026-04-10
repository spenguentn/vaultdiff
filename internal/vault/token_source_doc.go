// Package vault provides types and helpers for interacting with HashiCorp Vault.
//
// # TokenSource
//
// TokenSource captures how and when a Vault token was resolved, along with
// an optional TTL so callers can determine whether the token needs renewal.
//
// Supported source types:
//
//   - TokenSourceDirect  — token supplied inline in configuration
//   - TokenSourceEnv     — token read from an environment variable
//   - TokenSourceFile    — token read from a file on disk
//   - TokenSourceAppRole — token obtained via AppRole authentication
//
// Usage:
//
//	ts, err := vault.NewTokenSource(vault.TokenSourceEnv, token, 1*time.Hour)
//	if err != nil { ... }
//	if ts.IsExpired() { /* re-authenticate */ }
package vault
