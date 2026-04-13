// Package vault provides the SecretFreeze and SecretFreezeRegistry types for
// managing write-protection on individual secret paths.
//
// # Overview
//
// A SecretFreeze record marks a secret at a given mount and path as frozen,
// preventing modifications until explicitly unfrozen. Each record captures who
// froze the secret, the reason (manual, compliance, or incident), an optional
// human-readable note, and an optional expiry time after which the freeze is
// automatically considered inactive.
//
// # Usage
//
//	reg := vault.NewSecretFreezeRegistry()
//
//	err := reg.Freeze(vault.SecretFreeze{
//		Mount:    "secret",
//		Path:     "app/db-password",
//		FrozenBy: "ops-team",
//		Reason:   vault.FreezeReasonCompliance,
//		Note:     "PCI audit window",
//	})
//
//	if reg.IsFrozen("secret", "app/db-password") {
//		// block write operations
//	}
//
//	reg.Unfreeze("secret", "app/db-password")
package vault
