// Package vault provides rollback support for HashiCorp Vault KV v2 secrets.
//
// # Rollback
//
// A Rollbacker restores a secret to a previously stored version by reading
// the target version's data and writing it as a new version, which is the
// recommended pattern for KV v2 rollbacks.
//
// Basic usage:
//
//	rb := vault.NewRollbacker(client)
//	result := rb.Rollback(ctx, vault.RollbackRequest{
//		Mount:   "secret",
//		Path:    "myapp/database",
//		Version: 3,
//	})
//	if !result.IsSuccess() {
//		log.Fatalf("rollback failed: %v", result.Err)
//	}
//
// # RollbackPlan
//
// A RollbackPlan collects multiple rollback requests for batch execution.
// Each request is validated when added to the plan.
//
//	plan := vault.NewRollbackPlan()
//	_ = plan.Add(vault.RollbackRequest{Mount: "secret", Path: "app/db", Version: 2})
//	fmt.Print(plan.Describe())
package vault
