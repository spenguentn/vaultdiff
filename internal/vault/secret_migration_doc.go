// Package vault provides types and utilities for interacting with HashiCorp Vault.
//
// # Secret Migration
//
// SecretMigration tracks the movement of a secret from one mount/path to another,
// optionally across environments. Migrations progress through well-defined statuses:
//
//	 pending   → running → completed
//	                     ↘ failed
//
// Use SecretMigrationRegistry to submit, update, and query migrations. Only one
// active (non-terminal) migration is permitted per source path at a time.
//
// Example:
//
//	reg := vault.NewSecretMigrationRegistry()
//	result, err := reg.Submit(vault.SecretMigration{
//		SourceMount: "secret",
//		SourcePath:  "app/db-password",
//		DestMount:   "kv",
//		DestPath:    "prod/db-password",
//		InitiatedBy: "alice",
//		Status:      vault.MigrationPending,
//	})
package vault
