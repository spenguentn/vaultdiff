package vault

import "fmt"

// ValidCategories defines the allowed secret category values.
var ValidCategories = []string{
	"credentials",
	"certificates",
	"tokens",
	"keys",
	"config",
	"other",
}

// IsValidCategory reports whether the given category string is recognized.
func IsValidCategory(c string) bool {
	for _, v := range ValidCategories {
		if v == c {
			return true
		}
	}
	return false
}

// SecretCategory associates a secret path with a category label.
type SecretCategory struct {
	Mount    string
	Path     string
	Category string
	SetBy    string
}

// FullPath returns the canonical mount+path string.
func (s SecretCategory) FullPath() string {
	return s.Mount + "/" + s.Path
}

// Validate returns an error if the SecretCategory is missing required fields
// or contains an unrecognised category value.
func (s SecretCategory) Validate() error {
	if s.Mount == "" {
		return fmt.Errorf("secret category: mount is required")
	}
	if s.Path == "" {
		return fmt.Errorf("secret category: path is required")
	}
	if s.Category == "" {
		return fmt.Errorf("secret category: category is required")
	}
	if !IsValidCategory(s.Category) {
		return fmt.Errorf("secret category: unknown category %q", s.Category)
	}
	if s.SetBy == "" {
		return fmt.Errorf("secret category: set_by is required")
	}
	return nil
}

// categoryKey returns the registry lookup key for a SecretCategory.
func categoryKey(mount, path string) string {
	return mount + "::" + path
}

// SecretCategoryRegistry stores category assignments keyed by mount+path.
type SecretCategoryRegistry struct {
	entries map[string]SecretCategory
}

// NewSecretCategoryRegistry returns an initialised SecretCategoryRegistry.
func NewSecretCategoryRegistry() *SecretCategoryRegistry {
	return &SecretCategoryRegistry{entries: make(map[string]SecretCategory)}
}

// Set validates and stores a category assignment.
func (r *SecretCategoryRegistry) Set(sc SecretCategory) error {
	if err := sc.Validate(); err != nil {
		return err
	}
	r.entries[categoryKey(sc.Mount, sc.Path)] = sc
	return nil
}

// Get retrieves a category assignment. The second return value is false when
// no entry exists for the given mount and path.
func (r *SecretCategoryRegistry) Get(mount, path string) (SecretCategory, bool) {
	sc, ok := r.entries[categoryKey(mount, path)]
	return sc, ok
}

// Remove deletes the category assignment for the given mount and path.
func (r *SecretCategoryRegistry) Remove(mount, path string) {
	delete(r.entries, categoryKey(mount, path))
}

// Len returns the number of registered category entries.
func (r *SecretCategoryRegistry) Len() int {
	return len(r.entries)
}
