// Package vault provides types and utilities for interacting with
// HashiCorp Vault, including secret reading, metadata inspection,
// lease management, rollback planning, and policy modelling.
//
// # Policy
//
// The Policy and PolicyRule types represent Vault access control
// policies. A Policy is a named collection of PathRules, each
// granting one or more capabilities (read, list, create, update,
// delete, deny) on a Vault path.
//
// Use Validate to verify that a policy or rule is well-formed before
// submitting it to Vault or including it in a diff report.
package vault
