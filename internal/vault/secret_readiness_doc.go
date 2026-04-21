// Package vault provides types and registries for tracking secret readiness
// within a HashiCorp Vault environment.
//
// # Secret Readiness
//
// SecretReadiness records whether a secret at a given mount and path is
// considered ready for consumption by dependent services. Readiness is
// assessed externally (e.g. by a rotation pipeline or health-check job) and
// stored via SecretReadinessRegistry.
//
// Supported statuses:
//
//	- ReadinessStatusReady    – secret is valid and safe to use
//	- ReadinessStatusNotReady – secret exists but is not yet safe to use
//	- ReadinessStatusUnknown  – readiness has not been assessed
//
// # Usage
//
//	reg := vault.NewSecretReadinessRegistry()
//
//	err := reg.Set(vault.SecretReadiness{
//		Mount:  "secret",
//		Path:   "myapp/db-password",
//		Status: vault.ReadinessStatusReady,
//		Reason: "rotation complete",
//	})
//
//	if rec, ok := reg.Get("secret", "myapp/db-password"); ok && rec.IsReady() {
//		// proceed to consume the secret
//	}
package vault
