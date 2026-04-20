// Package vault provides types and registries for managing Vault secret
// lifecycle metadata.
//
// # Secret Disposition
//
// SecretDisposition records the intended end-of-life handling for a secret,
// capturing what action should be taken (delete, archive, rotate, or transfer),
// when it is scheduled, and who approved the decision.
//
// Use SecretDispositionRegistry to store, retrieve, and query dispositions.
// The Due method returns all records whose scheduled time has passed, enabling
// automated enforcement pipelines to act on overdue secrets.
//
// Example:
//
//	reg := vault.NewSecretDispositionRegistry()
//	err := reg.Set(&vault.SecretDisposition{
//		Mount:       "secret",
//		Path:        "app/legacy-token",
//		Action:      vault.DispositionDelete,
//		ScheduledAt: time.Now().Add(7 * 24 * time.Hour),
//		ApprovedBy:  "security-team",
//	})
package vault
