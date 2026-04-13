package vault

import (
	"errors"
	"fmt"
	"time"
)

// MigrationStatus represents the current state of a secret migration.
type MigrationStatus string

const (
	MigrationPending   MigrationStatus = "pending"
	MigrationRunning   MigrationStatus = "running"
	MigrationCompleted MigrationStatus = "completed"
	MigrationFailed    MigrationStatus = "failed"
)

// IsValidMigrationStatus returns true if the given status is known.
func IsValidMigrationStatus(s MigrationStatus) bool {
	switch s {
	case MigrationPending, MigrationRunning, MigrationCompleted, MigrationFailed:
		return true
	}
	return false
}

// SecretMigration describes a planned or completed migration of a secret
// from one mount/path to another, optionally across environments.
type SecretMigration struct {
	ID          string          `json:"id"`
	SourceMount string          `json:"source_mount"`
	SourcePath  string          `json:"source_path"`
	DestMount   string          `json:"dest_mount"`
	DestPath    string          `json:"dest_path"`
	Status      MigrationStatus `json:"status"`
	InitiatedBy string          `json:"initiated_by"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Note        string          `json:"note,omitempty"`
}

// FullSource returns the fully qualified source path.
func (m *SecretMigration) FullSource() string {
	return fmt.Sprintf("%s/%s", m.SourceMount, m.SourcePath)
}

// FullDest returns the fully qualified destination path.
func (m *SecretMigration) FullDest() string {
	return fmt.Sprintf("%s/%s", m.DestMount, m.DestPath)
}

// IsTerminal returns true if the migration has reached a final state.
func (m *SecretMigration) IsTerminal() bool {
	return m.Status == MigrationCompleted || m.Status == MigrationFailed
}

// Validate checks that the migration has all required fields.
func (m *SecretMigration) Validate() error {
	if m.SourceMount == "" {
		return errors.New("source_mount is required")
	}
	if m.SourcePath == "" {
		return errors.New("source_path is required")
	}
	if m.DestMount == "" {
		return errors.New("dest_mount is required")
	}
	if m.DestPath == "" {
		return errors.New("dest_path is required")
	}
	if m.InitiatedBy == "" {
		return errors.New("initiated_by is required")
	}
	if !IsValidMigrationStatus(m.Status) {
		return fmt.Errorf("unknown migration status: %q", m.Status)
	}
	return nil
}
