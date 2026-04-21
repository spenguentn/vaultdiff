package vault

import (
	"fmt"
	"time"
)

// SupersessionReason describes why a secret was superseded.
type SupersessionReason string

const (
	SupersessionReasonRotated   SupersessionReason = "rotated"
	SupersessionReasonMigrated  SupersessionReason = "migrated"
	SupersessionReasonDeprecated SupersessionReason = "deprecated"
	SupersessionReasonReplaced  SupersessionReason = "replaced"
)

// IsValidSupersessionReason returns true if the reason is a known value.
func IsValidSupersessionReason(r SupersessionReason) bool {
	switch r {
	case SupersessionReasonRotated, SupersessionReasonMigrated,
		SupersessionReasonDeprecated, SupersessionReasonReplaced:
		return true
	}
	return false
}

// SecretSupersession records that a secret at a given path has been
// superseded by another secret.
type SecretSupersession struct {
	Mount        string             `json:"mount"`
	Path         string             `json:"path"`
	SupersededBy string             `json:"superseded_by"`
	Reason       SupersessionReason `json:"reason"`
	SupersededAt time.Time          `json:"superseded_at"`
	Actor        string             `json:"actor"`
}

// FullPath returns the canonical mount+path string.
func (s SecretSupersession) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Mount, s.Path)
}

// Validate returns an error if the supersession record is incomplete.
func (s SecretSupersession) Validate() error {
	if s.Mount == "" {
		return fmt.Errorf("supersession: mount is required")
	}
	if s.Path == "" {
		return fmt.Errorf("supersession: path is required")
	}
	if s.SupersededBy == "" {
		return fmt.Errorf("supersession: superseded_by is required")
	}
	if !IsValidSupersessionReason(s.Reason) {
		return fmt.Errorf("supersession: unknown reason %q", s.Reason)
	}
	if s.Actor == "" {
		return fmt.Errorf("supersession: actor is required")
	}
	return nil
}
