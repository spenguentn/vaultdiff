// Package vault provides types and utilities for interacting with
// HashiCorp Vault, including secret reading, diffing, caching, leasing,
// and schema validation.
//
// # Secret Schema
//
// SecretSchema describes the expected structure and types of keys within a
// Vault secret. It can be used to validate secret data before writing or
// after reading, ensuring that required fields are present and values conform
// to expected formats.
//
// Supported field types:
//
//   - FieldTypeString  – any non-empty string value (no format check)
//   - FieldTypeNumeric – integer or decimal number
//   - FieldTypeBoolean – must be the literal string "true" or "false"
//   - FieldTypePattern – value must match the provided regular expression
//
// Example:
//
//	schema := vault.SecretSchema{
//		Mount: "secret",
//		Path:  "app/config",
//		Fields: []vault.FieldSchema{
//			{Key: "db_host", Type: vault.FieldTypeString,  Required: true},
//			{Key: "port",    Type: vault.FieldTypeNumeric, Required: false},
//		},
//	}
//
//	if err := schema.Validate(); err != nil {
//		log.Fatal(err)
//	}
//
//	violations := schema.ValidateSecretData(data)
package vault
