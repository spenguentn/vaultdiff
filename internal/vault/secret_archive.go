package vault

import (
	"fmt"
	"time"
)

// ArchiveReason describes why a secret was archived.
type ArchiveReason string

const (
	ArchiveReasonDeprecated ArchiveReason = "deprecated"
	ArchiveReasonRotated    ArchiveReason = "rotated"
	ArchiveReasonMigrated   ArchiveReason = "migrated"
	ArchiveReasonManual     ArchiveReason = "manual"
)

// SecretArchiveEntry records a secret that has been archived.
type SecretArchiveEntry struct {
	Mount      string            `json:"mount"`
	Path       string            `json:"path"`
	Version    int               `json:"version"`
	Reason     ArchiveReason     `json:"reason"`
	ArchivedBy string            `json:"archived_by"`
	ArchivedAt time.Time         `json:"archived_at"`
	Labels     map[string]string `json:"labels,omitempty"`
}

// Validate checks that the archive entry has all required fields.
func (e *SecretArchiveEntry) Validate() error {
	if e.Mount == "" {
		return fmt.Errorf("archive entry: mount is required")
	}
	if e.Path == "" {
		return fmt.Errorf("archive entry: path is required")
	}
	if e.Version < 0 {
		return fmt.Errorf("archive entry: version must be non-negative")
	}
	if e.ArchivedBy == "" {
		return fmt.Errorf("archive entry: archived_by is required")
	}
	if e.Reason == "" {
		return fmt.Errorf("archive entry: reason is required")
	}
	return nil
}

// FullPath returns the combined mount and path.
func (e *SecretArchiveEntry) FullPath() string {
	return e.Mount + "/" + e.Path
}

// IsReasonValid reports whether the reason is a known archive reason.
func IsReasonValid(r ArchiveReason) bool {
	switch r {
	case ArchiveReasonDeprecated, ArchiveReasonRotated, ArchiveReasonMigrated, ArchiveReasonManual:
		return true
	}
	return false
}
