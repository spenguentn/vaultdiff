// Package vault provides types and registries for managing HashiCorp Vault
// secrets, metadata, and governance records.
//
// # Secret Covenants
//
// A SecretCovenant captures a governance agreement associated with a specific
// secret path.  Covenants describe who owns a secret, under what access
// model it is shared, and optionally when that agreement expires.
//
// Supported covenant types:
//
//   - shared    – multiple consumers may read the secret.
//   - exclusive – a single consumer owns read access.
//   - read_only – the secret must not be modified by consumers.
//
// Usage:
//
//	reg := vault.NewSecretCovenantRegistry()
//
//	err := reg.Set(&vault.SecretCovenant{
//		Mount: "secret",
//		Path:  "app/api-key",
//		Type:  vault.CovenantTypeExclusive,
//		Owner: "team-backend",
//	})
//
//	covenant, err := reg.Get("secret", "app/api-key")
package vault
