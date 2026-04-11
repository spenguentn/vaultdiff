package vault

import (
	"testing"
)

var baseSchema = SecretSchema{
	Mount: "secret",
	Path:  "app/config",
	Fields: []FieldSchema{
		{Key: "db_host", Type: FieldTypeString, Required: true},
		{Key: "port", Type: FieldTypeNumeric, Required: false},
		{Key: "debug", Type: FieldTypeBoolean, Required: false},
		{Key: "api_key", Type: FieldTypePattern, Required: true, Pattern: `^[A-Z0-9]{16}$`},
	},
}

func TestFieldSchema_Validate_Valid(t *testing.T) {
	f := FieldSchema{Key: "host", Type: FieldTypeString}
	if err := f.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFieldSchema_Validate_EmptyKey(t *testing.T) {
	f := FieldSchema{Key: "", Type: FieldTypeString}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestFieldSchema_Validate_PatternMissingRegex(t *testing.T) {
	f := FieldSchema{Key: "token", Type: FieldTypePattern, Pattern: ""}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for pattern type with no pattern")
	}
}

func TestFieldSchema_Validate_InvalidRegex(t *testing.T) {
	f := FieldSchema{Key: "token", Type: FieldTypePattern, Pattern: `[invalid`}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestSecretSchema_Validate_Valid(t *testing.T) {
	if err := baseSchema.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretSchema_Validate_MissingMount(t *testing.T) {
	s := baseSchema
	s.Mount = ""
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestSecretSchema_Validate_MissingPath(t *testing.T) {
	s := baseSchema
	s.Path = ""
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestValidateSecretData_Valid(t *testing.T) {
	data := map[string]string{
		"db_host": "localhost",
		"port":    "5432",
		"debug":   "false",
		"api_key": "ABCD1234EFGH5678",
	}
	errs := baseSchema.ValidateSecretData(data)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}
}

func TestValidateSecretData_MissingRequired(t *testing.T) {
	data := map[string]string{
		"port": "5432",
	}
	errs := baseSchema.ValidateSecretData(data)
	if len(errs) == 0 {
		t.Fatal("expected errors for missing required fields")
	}
}

func TestValidateSecretData_InvalidBoolean(t *testing.T) {
	data := map[string]string{
		"db_host": "localhost",
		"api_key": "ABCD1234EFGH5678",
		"debug":   "yes",
	}
	errs := baseSchema.ValidateSecretData(data)
	if len(errs) == 0 {
		t.Fatal("expected error for invalid boolean value")
	}
}

func TestValidateSecretData_InvalidNumeric(t *testing.T) {
	data := map[string]string{
		"db_host": "localhost",
		"api_key": "ABCD1234EFGH5678",
		"port":    "not-a-number",
	}
	errs := baseSchema.ValidateSecretData(data)
	if len(errs) == 0 {
		t.Fatal("expected error for invalid numeric value")
	}
}

func TestValidateSecretData_PatternMismatch(t *testing.T) {
	data := map[string]string{
		"db_host": "localhost",
		"api_key": "short",
	}
	errs := baseSchema.ValidateSecretData(data)
	if len(errs) == 0 {
		t.Fatal("expected error for pattern mismatch")
	}
}
