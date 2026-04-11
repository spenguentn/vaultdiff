package vault

import (
	"errors"
	"fmt"
	"time"
)

// LineageEntry represents a single point in a secret's version history.
type LineageEntry struct {
	Mount     string
	Path      string
	Version   int
	CreatedAt time.Time
	DeletedAt *time.Time
	CreatedBy string
}

// FullPath returns the canonical mount+path string for the entry.
func (e LineageEntry) FullPath() string {
	return fmt.Sprintf("%s/%s", e.Mount, e.Path)
}

// IsDeleted reports whether this version has been soft-deleted.
func (e LineageEntry) IsDeleted() bool {
	return e.DeletedAt != nil
}

// Validate checks that the entry has the minimum required fields.
func (e LineageEntry) Validate() error {
	if e.Mount == "" {
		return errors.New("lineage entry: mount is required")
	}
	if e.Path == "" {
		return errors.New("lineage entry: path is required")
	}
	if e.Version < 1 {
		return errors.New("lineage entry: version must be >= 1")
	}
	return nil
}

// SecretLineage holds the ordered version history for a single secret.
type SecretLineage struct {
	Mount   string
	Path    string
	entries []LineageEntry
}

// NewSecretLineage creates a SecretLineage for the given mount and path.
func NewSecretLineage(mount, path string) *SecretLineage {
	return &SecretLineage{Mount: mount, Path: path}
}

// Add appends a validated entry to the lineage.
func (l *SecretLineage) Add(e LineageEntry) error {
	if err := e.Validate(); err != nil {
		return err
	}
	l.entries = append(l.entries, e)
	return nil
}

// Entries returns a copy of all lineage entries.
func (l *SecretLineage) Entries() []LineageEntry {
	out := make([]LineageEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Latest returns the most recently added entry, if any.
func (l *SecretLineage) Latest() (LineageEntry, bool) {
	if len(l.entries) == 0 {
		return LineageEntry{}, false
	}
	return l.entries[len(l.entries)-1], true
}

// Len returns the number of entries in the lineage.
func (l *SecretLineage) Len() int {
	return len(l.entries)
}
