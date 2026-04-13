// Package vault provides the SecretLabel and SecretLabelRegistry types for
// attaching arbitrary key-value metadata to Vault secrets.
//
// Labels are lightweight annotations intended for organizational purposes such
// as team ownership, environment tagging, and cost-centre attribution. Unlike
// annotations (which carry richer structured data), labels are simple strings.
//
// Usage:
//
//	reg := vault.NewSecretLabelRegistry()
//	err := reg.Set(vault.SecretLabel{
//		Mount:     "secret",
//		Path:      "app/db",
//		Key:       "team",
//		Value:     "platform",
//		CreatedBy: "alice",
//	})
//
//	labels := reg.List("secret", "app/db")
package vault
