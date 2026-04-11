// Package vault provides the SecretApprovalRegistry for managing change-approval
// workflows on Vault secrets.
//
// # Overview
//
// The approval workflow allows operators to gate secret mutations behind a
// review step. A change author submits an ApprovalRequest; a reviewer then
// approves or rejects it before the mutation is applied.
//
// # Lifecycle
//
//	┌──────────┐  Submit   ┌─────────┐  Review(true)  ┌──────────┐
//	│  (none)  │ ────────► │ pending │ ─────────────► │ approved │
//	└──────────┘           └─────────┘                └──────────┘
//	                            │  Review(false)            │ Revoke
//	                            ▼                           ▼
//	                       ┌──────────┐             ┌─────────────┐
//	                       │ rejected │             │   revoked   │
//	                       └──────────┘             └─────────────┘
//
// # Usage
//
//	reg := vault.NewSecretApprovalRegistry()
//
//	req := &vault.ApprovalRequest{
//	    Mount:       "secret",
//	    Path:        "myapp/db",
//	    RequestedBy: "alice",
//	    Reason:      "quarterly rotation",
//	}
//	reg.Submit(req)
//	reg.Review("secret", "myapp/db", "bob", true)
package vault
