package vault

import (
	"fmt"
	"regexp"
)

// ValidationResult holds the outcome of validating a secret against a schema.
type ValidationResult struct {
	Mount   string
	Path    string
	Passed  bool
	Errors  []string
}

// IsValid returns true when there are no validation errors.
func (r ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// SecretValidator validates secret data against a SecretSchema.
type SecretValidator struct {
	schema SecretSchema
}

// NewSecretValidator creates a SecretValidator for the given schema.
// It returns an error if the schema is invalid.
func NewSecretValidator(schema SecretSchema) (*SecretValidator, error) {
	if err := schema.Validate(); err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}
	return &SecretValidator{schema: schema}, nil
}

// Validate checks the provided secret data map against the schema rules.
// It returns a ValidationResult describing any failures.
func (v *SecretValidator) Validate(mount, path string, data map[string]string) ValidationResult {
	result := ValidationResult{
		Mount:  mount,
		Path:   path,
		Passed: true,
	}

	for _, field := range v.schema.Fields {
		val, exists := data[field.Key]

		if field.Required && !exists {
			result.Errors = append(result.Errors, fmt.Sprintf("required field %q is missing", field.Key))
			continue
		}

		if !exists {
			continue
		}

		if field.Pattern != "" {
			matched, err := regexp.MatchString(field.Pattern, val)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("field %q has invalid pattern: %v", field.Key, err))
				continue
			}
			if !matched {
				result.Errors = append(result.Errors,
					fmt.Sprintf("field %q value does not match required pattern %q", field.Key, field.Pattern))
			}
		}
	}

	if len(result.Errors) > 0 {
		result.Passed = false
	}
	return result
}
