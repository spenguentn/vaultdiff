// Package vault provides the SecretTrust type and SecretTrustRegistry for
// tracking the trust level assigned to secrets stored in HashiCorp Vault.
//
// Trust levels range from "untrusted" through "verified" and are recorded
// alongside the identity of the operator who made the assignment. The
// registry is safe for concurrent use.
//
// Example usage:
//
//	reg := vault.NewSecretTrustRegistry()
//	err := reg.Set(&vault.SecretTrust{
//		Mount:      "secret",
//		Path:       "services/api-key",
//		Level:      vault.TrustLevelHigh,
//		AssignedBy: "platform-team",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
package vault
