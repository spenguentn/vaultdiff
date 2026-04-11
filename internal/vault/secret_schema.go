package vault

import (
	"errors"
	"fmt"
	"regexp"
)

// FieldType represents the expected type of a secret field.
type FieldType string

const (
	FieldTypeString FieldType = "string"
	FieldTypeNumeric FieldType = "numeric"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypePattern FieldType = "pattern"
)

// FieldSchema defines validation rules for a single secret key.
type FieldSchema struct {
	Key      string
	Type     FieldType
	Required bool
	Pattern  string // used when Type == FieldTypePattern
}

// Validate checks that the FieldSchema is well-formed.
func (f FieldSchema) Validate() error {
	if f.Key == "" {
		return errors.New("field schema: key must not be empty")
	}
	if f.Type == "" {
		return errors.New("field schema: type must not be empty")
	}
	if f.Type == FieldTypePattern && f.Pattern == "" {
		return fmt.Errorf("field schema: key %q has type pattern but no pattern set", f.Key)
	}
	if f.Type == FieldTypePattern {
		if _, err := regexp.Compile(f.Pattern); err != nil {
			return fmt.Errorf("field schema: key %q has invalid pattern: %w", f.Key, err)
		}
	}
	return nil
}

// SecretSchema holds the full schema for a secret at a given path.
type SecretSchema struct {
	Mount  string
	Path   string
	Fields []FieldSchema
}

// Validate checks the schema and all its fields.
func (s SecretSchema) Validate() error {
	if s.Mount == "" {
		return errors.New("secret schema: mount must not be empty")
	}
	if s.Path == "" {
		return errors.New("secret schema: path must not be empty")
	}
	for _, f := range s.Fields {
		if err := f.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ValidateSecretData checks a map of secret key/value pairs against the schema.
// It returns a slice of validation error strings (empty means valid).
func (s SecretSchema) ValidateSecretData(data map[string]string) []string {
	var errs []string
	for _, field := range s.Fields {
		val, ok := data[field.Key]
		if !ok {
			if field.Required {
				errs = append(errs, fmt.Sprintf("missing required key: %q", field.Key))
			}
			continue
		}
		switch field.Type {
		case FieldTypeBoolean:
			if val != "true" && val != "false" {
				errs = append(errs, fmt.Sprintf("key %q must be 'true' or 'false', got %q", field.Key, val))
			}
		case FieldTypeNumeric:
			if !regexp.MustCompile(`^-?[0-9]+(\.[0-9]+)?$`).MatchString(val) {
				errs = append(errs, fmt.Sprintf("key %q must be numeric, got %q", field.Key, val))
			}
		case FieldTypePattern:
			if !regexp.MustCompile(field.Pattern).MatchString(val) {
				errs = append(errs, fmt.Sprintf("key %q does not match pattern %q", field.Key, field.Pattern))
			}
		}
	}
	return errs
}
