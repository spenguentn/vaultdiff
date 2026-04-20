// Package vault provides the SecretWorkflow type and SecretWorkflowRegistry
// for tracking the end-to-end lifecycle workflow of a Vault secret.
//
// A workflow moves a secret through a series of well-defined stages:
//
//	 draft → pending → approved → deployed → retired
//	                 ↘ rejected
//
// Each transition is recorded with the actor who performed it and an optional
// comment, giving an auditable trail of every stage change.
//
// Usage:
//
//	reg := vault.NewSecretWorkflowRegistry()
//
//	_ = reg.Set(vault.SecretWorkflow{
//		Mount:       "secret",
//		Path:        "app/api-key",
//		Stage:       vault.WorkflowStageDraft,
//		InitiatedBy: "alice",
//	})
//
//	_ = reg.Advance("secret", "app/api-key", vault.WorkflowStageApproved, "bob", "approved")
package vault
