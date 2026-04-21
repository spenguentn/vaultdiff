// Package vault provides the SecretRevocationRegistry for tracking
// revoked secrets within a Vault deployment.
//
// A RevocationRecord captures who revoked a secret, why it was revoked,
// when the revocation occurred, and an optional expiry after which the
// revocation record itself is considered stale.
//
// Supported revocation reasons:
//
//	"compromised"       — secret was exposed or leaked
//	"expired"           — secret passed its natural TTL
//	"rotated"           — secret was superseded by a new version
//	"decommissioned"    — the service using the secret was retired
//	"policy_violation"  — secret violated a security or compliance policy
//
// Usage:
//
//	reg := vault.NewSecretRevocationRegistry()
//	err := reg.Revoke(vault.RevocationRecord{
//		Mount:     "secret",
//		Path:      "app/db-password",
//		RevokedBy: "security-team",
//		Reason:    "compromised",
//		RevokedAt: time.Now(),
//	})
package vault
