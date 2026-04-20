package vault

import (
	"fmt"
	"time"
)

// ObsolescenceReason describes why a secret is considered obsolete.
type ObsolescenceReason string

const (
	ObsolescenceReasonSuperseded ObsolescenceReason = "superseded"
	ObsolescenceReasonDeprecated ObsolescenceReason = "deprecated"
	ObsolescenceReasonUnused     ObsolescenceReason = "unused"
	ObsolescenceReasonExpired    ObsolescenceReason = "expired"
)

// IsValidObsolescenceReason returns true if the given reason is recognised.
func IsValidObsolescenceReason(r ObsolescenceReason) bool {
	switch r {
	case ObsolescenceReasonSuperseded,
		ObsolescenceReasonDeprecated,
		ObsolescenceReasonUnused,
		ObsolescenceReasonExpired:
		return true
	}
	return false
}

// SecretObsolescence records that a secret has been marked obsolete.
type SecretObsolescence struct {
	Mount       string             `json:"mount"`
	Path        string             `json:"path"`
	Reason      ObsolescenceReason `json:"reason"`
	MarkedBy    string             `json:"marked_by"`
	MarkedAt    time.Time          `json:"marked_at"`
	ReplacedBy  string             `json:"replaced_by,omitempty"`
	ScheduledAt *time.Time         `json:"scheduled_removal_at,omitempty"`
}

// FullPath returns the canonical mount+path identifier.
func (o *SecretObsolescence) FullPath() string {
	return fmt.Sprintf("%s/%s", o.Mount, o.Path)
}

// IsDue reports whether the scheduled removal time has passed.
func (o *SecretObsolescence) IsDue() bool {
	if o.ScheduledAt == nil {
		return false
	}
	return time.Now().After(*o.ScheduledAt)
}

// Validate checks that the record is well-formed.
func (o *SecretObsolescence) Validate() error {
	if o.Mount == "" {
		return fmt.Errorf("obsolescence: mount is required")
	}
	if o.Path == "" {
		return fmt.Errorf("obsolescence: path is required")
	}
	if !IsValidObsolescenceReason(o.Reason) {
		return fmt.Errorf("obsolescence: unknown reason %q", o.Reason)
	}
	if o.MarkedBy == "" {
		return fmt.Errorf("obsolescence: marked_by is required")
	}
	return nil
}
