package vault

import (
	"fmt"
	"strings"
	"time"
)

// SecretLabel represents a key-value label attached to a secret for
// organizational and filtering purposes.
type SecretLabel struct {
	Mount     string
	Path      string
	Key       string
	Value     string
	CreatedBy string
	CreatedAt time.Time
}

// FullPath returns the canonical vault path for the labelled secret.
func (l SecretLabel) FullPath() string {
	return fmt.Sprintf("%s/%s", strings.Trim(l.Mount, "/"), strings.Trim(l.Path, "/"))
}

// Validate checks that required fields are present.
func (l SecretLabel) Validate() error {
	if l.Mount == "" {
		return fmt.Errorf("secret label: mount is required")
	}
	if l.Path == "" {
		return fmt.Errorf("secret label: path is required")
	}
	if l.Key == "" {
		return fmt.Errorf("secret label: key is required")
	}
	if l.CreatedBy == "" {
		return fmt.Errorf("secret label: created_by is required")
	}
	return nil
}

func labelKey(mount, path, key string) string {
	return fmt.Sprintf("%s/%s#%s", strings.Trim(mount, "/"), strings.Trim(path, "/"), key)
}
