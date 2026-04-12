package vault

import (
	"testing"
)

func baseValidatorSchema() SecretSchema {
	return SecretSchema{
		Mount: "secret",
		Path:  "app/config",
		Fields: []FieldSchema{
			{Key: "api_key", Required: true, Pattern: `^[A-Za-z0-9]{16,}$`},
			{Key: "env", Required: true},
			{Key: "debug", Required: false},
		},
	}
}

func TestNewSecretValidator_Valid(t *testing.T) {
	v, err := NewSecretValidator(baseValidatorSchema())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected non-nil validator")
	}
}

func TestNewSecretValidator_InvalidSchema(t *testing.T) {
	bad := SecretSchema{} // missing mount/path/fields
	_, err := NewSecretValidator(bad)
	if err == nil {
		t.Fatal("expected error for invalid schema")
	}
}

func TestValidate_AllFieldsValid(t *testing.T) {
	v, _ := NewSecretValidator(baseValidatorSchema())
	data := map[string]string{
		"api_key": "ABCDEFGH12345678",
		"env":     "production",
	}
	res := v.Validate("secret", "app/config", data)
	if !res.IsValid() {
		t.Fatalf("expected valid, got errors: %v", res.Errors)
	}
}

func TestValidate_MissingRequiredField(t *testing.T) {
	v, _ := NewSecretValidator(baseValidatorSchema())
	data := map[string]string{
		"env": "staging",
	}
	res := v.Validate("secret", "app/config", data)
	if res.IsValid() {
		t.Fatal("expected invalid due to missing api_key")
	}
	if len(res.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(res.Errors))
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	v, _ := NewSecretValidator(baseValidatorSchema())
	data := map[string]string{
		"api_key": "short",
		"env":     "dev",
	}
	res := v.Validate("secret", "app/config", data)
	if res.IsValid() {
		t.Fatal("expected invalid due to pattern mismatch")
	}
}

func TestValidate_OptionalFieldSkipped(t *testing.T) {
	v, _ := NewSecretValidator(baseValidatorSchema())
	data := map[string]string{
		"api_key": "ABCDEFGH12345678",
		"env":     "test",
		// debug omitted — optional
	}
	res := v.Validate("secret", "app/config", data)
	if !res.IsValid() {
		t.Fatalf("expected valid, got: %v", res.Errors)
	}
}

func TestValidationResult_Fields(t *testing.T) {
	v, _ := NewSecretValidator(baseValidatorSchema())
	res := v.Validate("secret", "app/config", map[string]string{})
	if res.Mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", res.Mount)
	}
	if res.Path != "app/config" {
		t.Errorf("expected path 'app/config', got %q", res.Path)
	}
}
