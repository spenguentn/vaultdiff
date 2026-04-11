// Package vault provides types and utilities for interacting with
// HashiCorp Vault secrets, including archiving capabilities.
//
// # Secret Archive
//
// SecretArchiveEntry records metadata about a secret that has been
// intentionally retired from active use. Each entry captures:
//
//   - The Vault mount and path of the secret
//   - The specific version that was archived
//   - The reason for archival (deprecated, rotated, migrated, manual)
//   - Who performed the archival and when
//
// SecretArchiveRegistry provides a thread-safe, in-memory store for
// archive entries. It supports archiving, retrieval, removal, and
// enumeration of entries.
//
// Example usage:
//
//	reg := vault.NewSecretArchiveRegistry()
//	err := reg.Archive(&vault.SecretArchiveEntry{
//		Mount:      "secret",
//		Path:       "myapp/db-password",
//		Version:    5,
//		Reason:     vault.ArchiveReasonRotated,
//		ArchivedBy: "ci-pipeline",
//	})
package vault
