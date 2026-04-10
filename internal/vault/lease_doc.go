// Package vault provides utilities for interacting with HashiCorp Vault,
// including secret reading, KV configuration, caching, retry logic,
// health checks, and lease lifecycle management.
//
// # Lease Management
//
// LeaseInfo represents the metadata associated with a dynamic Vault secret
// lease, including its ID, duration, renewability, and issue time.
//
// ParseLease constructs a LeaseInfo from raw Vault API response fields.
//
// LeaseTracker provides a thread-safe registry for tracking active leases.
// Use Track to register a lease, Get to retrieve it, Revoke to remove it,
// and Expired to enumerate leases past their TTL.
package vault
