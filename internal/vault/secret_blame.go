package vault

import (
	"fmt"
	"time"
)

// BlameEntry records who last modified a secret and when.
type BlameEntry struct {
	Mount     string
	Path      string
	Version   int
	ChangedBy string
	ChangedAt time.Time
	Operation string // e.g. "write", "delete", "restore"
}

// FullPath returns the canonical mount+path string.
func (b BlameEntry) FullPath() string {
	return fmt.Sprintf("%s/%s", b.Mount, b.Path)
}

// Validate checks that the BlameEntry has required fields.
func (b BlameEntry) Validate() error {
	if b.Mount == "" {
		return fmt.Errorf("blame entry: mount is required")
	}
	if b.Path == "" {
		return fmt.Errorf("blame entry: path is required")
	}
	if b.ChangedBy == "" {
		return fmt.Errorf("blame entry: changed_by is required")
	}
	if b.ChangedAt.IsZero() {
		return fmt.Errorf("blame entry: changed_at is required")
	}
	if b.Version < 1 {
		return fmt.Errorf("blame entry: version must be >= 1")
	}
	return nil
}

// IsValidOperation returns true if the operation string is a known type.
func IsValidOperation(op string) bool {
	switch op {
	case "write", "delete", "restore":
		return true
	}
	return false
}
