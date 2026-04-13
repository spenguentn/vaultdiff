package vault

import (
	"errors"
	"fmt"
	"time"
)

// SecretBookmark represents a named shortcut to a frequently accessed secret.
type SecretBookmark struct {
	Mount     string    `json:"mount"`
	Path      string    `json:"path"`
	Alias     string    `json:"alias"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	Note      string    `json:"note,omitempty"`
}

// FullPath returns the canonical mount+path string.
func (b SecretBookmark) FullPath() string {
	return fmt.Sprintf("%s/%s", b.Mount, b.Path)
}

// Validate returns an error if the bookmark is missing required fields.
func (b SecretBookmark) Validate() error {
	if b.Mount == "" {
		return errors.New("bookmark: mount is required")
	}
	if b.Path == "" {
		return errors.New("bookmark: path is required")
	}
	if b.Alias == "" {
		return errors.New("bookmark: alias is required")
	}
	if b.CreatedBy == "" {
		return errors.New("bookmark: created_by is required")
	}
	return nil
}
