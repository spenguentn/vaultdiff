// Package vault provides clients and helpers for interacting with
// HashiCorp Vault KV v2 secrets engines.
//
// It includes:
//   - Client construction and configuration (client.go)
//   - Secret reading and version parsing (secrets.go)
//   - Version listing (versions.go)
//   - Secret metadata retrieval (metadata.go)
//   - Environment definition and validation (env.go)
//   - Environment pair management for cross-environment diffs (env_pair.go)
package vault
