package vault

import "fmt"

// VisibilityLevel represents who can see a secret.
type VisibilityLevel string

const (
	VisibilityPublic   VisibilityLevel = "public"
	VisibilityInternal VisibilityLevel = "internal"
	VisibilityPrivate  VisibilityLevel = "private"
	VisibilityRestricted VisibilityLevel = "restricted"
)

// IsValidVisibilityLevel returns true if the given level is known.
func IsValidVisibilityLevel(l VisibilityLevel) bool {
	switch l {
	case VisibilityPublic, VisibilityInternal, VisibilityPrivate, VisibilityRestricted:
		return true
	}
	return false
}

// SecretVisibility records the visibility level assigned to a secret.
type SecretVisibility struct {
	Mount     string          `json:"mount"`
	Path      string          `json:"path"`
	Level     VisibilityLevel `json:"level"`
	SetBy     string          `json:"set_by"`
}

// FullPath returns the canonical mount+path string.
func (v SecretVisibility) FullPath() string {
	return fmt.Sprintf("%s/%s", v.Mount, v.Path)
}

// Validate returns an error if the record is incomplete or invalid.
func (v SecretVisibility) Validate() error {
	if v.Mount == "" {
		return fmt.Errorf("visibility: mount is required")
	}
	if v.Path == "" {
		return fmt.Errorf("visibility: path is required")
	}
	if !IsValidVisibilityLevel(v.Level) {
		return fmt.Errorf("visibility: unknown level %q", v.Level)
	}
	if v.SetBy == "" {
		return fmt.Errorf("visibility: set_by is required")
	}
	return nil
}
