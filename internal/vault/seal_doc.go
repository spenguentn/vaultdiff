// Package vault provides utilities for interacting with HashiCorp Vault,
// including secret reading, metadata inspection, health checking, and
// seal status management.
//
// # Seal Status
//
// The seal sub-feature exposes two main types:
//
//   - [SealInfo]: a value type representing a point-in-time snapshot of the
//     Vault seal state, including whether it is initialized, sealed, the
//     unseal progress, key threshold, total shares, and the running version.
//
//   - [SealChecker]: a thin wrapper around a [vault/api.Client] that fetches
//     the seal status from the /v1/sys/seal-status endpoint and returns a
//     parsed [SealInfo]. It also provides [SealChecker.MustBeUnsealed] as a
//     convenience guard used by commands that require an operational Vault.
//
// Example usage:
//
//	checker, err := vault.NewSealChecker(client, 5*time.Second)
//	if err != nil { ... }
//	if err := checker.MustBeUnsealed(ctx); err != nil { ... }
package vault
