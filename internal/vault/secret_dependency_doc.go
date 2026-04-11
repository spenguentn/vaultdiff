// Package vault provides the SecretDependencyRegistry for tracking
// directional relationships between secrets across mounts.
//
// # Overview
//
// A DependencyLink records that one secret (source) depends on another
// (target). This is useful for auditing blast-radius when a secret changes:
// any source that depends on the changed target may need to be rotated or
// reviewed.
//
// # Usage
//
//	reg := vault.NewSecretDependencyRegistry()
//
//	err := reg.Add(vault.DependencyLink{
//		SourceMount: "secret",
//		SourcePath:  "app/api",
//		TargetMount: "secret",
//		TargetPath:  "shared/db",
//		AddedBy:     "ops-team",
//		Note:        "api service reads db credentials at startup",
//	})
//
//	deps := reg.GetDependencies("secret", "app/api")
//	for _, d := range deps {
//		fmt.Println(d.FullTarget())
//	}
package vault
