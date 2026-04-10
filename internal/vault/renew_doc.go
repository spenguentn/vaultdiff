// Package vault provides primitives for interacting with HashiCorp Vault,
// including secret reads, metadata queries, lease management, and renewal.
//
// # Lease Renewal
//
// The Renewer type wraps the sys/leases/renew endpoint and provides a
// structured result type so callers can handle errors without inspecting
// raw HTTP responses.
//
// Basic usage:
//
//	 renewer := vault.NewRenewer(logicalClient)
//	 result := renewer.Renew(vault.RenewRequest{
//	     LeaseID:   "database/creds/my-role/abc123",
//	     Increment: 1 * time.Hour,
//	 })
//	 if !result.IsSuccess() {
//	     log.Printf("renewal failed: %v", result.Err)
//	 }
//
// The Increment field is advisory; Vault enforces its own maximum TTL and
// may return a shorter duration than requested.
package vault
