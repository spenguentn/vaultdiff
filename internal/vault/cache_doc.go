// Package vault provides clients and utilities for interacting with
// HashiCorp Vault, including KV secret reading, metadata inspection,
// version listing, and environment-aware configuration.
//
// The SecretCache type offers an optional in-memory caching layer for
// secret data retrieved from Vault. Entries are keyed by secret path
// and can be configured with a TTL to limit stale reads. A zero TTL
// means entries are retained indefinitely until explicitly invalidated
// or flushed.
//
// Cache usage is opt-in; callers that do not require caching can read
// secrets directly via the KV or multi-reader helpers.
package vault
