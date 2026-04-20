package vault

import (
	"fmt"
	"time"
)

// DispositionAction represents what should happen to a secret at end-of-life.
type DispositionAction string

const (
	DispositionDelete   DispositionAction = "delete"
	DispositionArchive  DispositionAction = "archive"
	DispositionRotate   DispositionAction = "rotate"
	DispositionTransfer DispositionAction = "transfer"
)

// IsValidDispositionAction returns true if the action is a known disposition action.
func IsValidDispositionAction(a DispositionAction) bool {
	switch a {
	case DispositionDelete, DispositionArchive, DispositionRotate, DispositionTransfer:
		return true
	}
	return false
}

// SecretDisposition defines the intended end-of-life handling for a secret.
type SecretDisposition struct {
	Mount       string            `json:"mount"`
	Path        string            `json:"path"`
	Action      DispositionAction `json:"action"`
	ScheduledAt time.Time         `json:"scheduled_at"`
	ApprovedBy  string            `json:"approved_by"`
	Notes       string            `json:"notes,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

// FullPath returns the canonical mount+path identifier.
func (d *SecretDisposition) FullPath() string {
	return fmt.Sprintf("%s/%s", d.Mount, d.Path)
}

// IsDue returns true if the disposition is scheduled at or before the given time.
func (d *SecretDisposition) IsDue(now time.Time) bool {
	return !d.ScheduledAt.IsZero() && !now.Before(d.ScheduledAt)
}

// Validate checks that all required fields are present and valid.
func (d *SecretDisposition) Validate() error {
	if d.Mount == "" {
		return fmt.Errorf("disposition: mount is required")
	}
	if d.Path == "" {
		return fmt.Errorf("disposition: path is required")
	}
	if !IsValidDispositionAction(d.Action) {
		return fmt.Errorf("disposition: invalid action %q", d.Action)
	}
	if d.ScheduledAt.IsZero() {
		return fmt.Errorf("disposition: scheduled_at is required")
	}
	if d.ApprovedBy == "" {
		return fmt.Errorf("disposition: approved_by is required")
	}
	return nil
}
