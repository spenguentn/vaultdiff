package vault

import (
	"errors"
	"fmt"
	"time"
)

// AuditEventKind classifies the type of auditable action on a secret.
type AuditEventKind string

const (
	AuditEventRead    AuditEventKind = "read"
	AuditEventWrite   AuditEventKind = "write"
	AuditEventDelete  AuditEventKind = "delete"
	AuditEventPromote AuditEventKind = "promote"
	AuditEventRollback AuditEventKind = "rollback"
)

// IsValidAuditEvent reports whether kind is a known audit event kind.
func IsValidAuditEvent(kind AuditEventKind) bool {
	switch kind {
	case AuditEventRead, AuditEventWrite, AuditEventDelete, AuditEventPromote, AuditEventRollback:
		return true
	}
	return false
}

// AuditTrailEntry records a single auditable action on a secret.
type AuditTrailEntry struct {
	Mount     string         `json:"mount"`
	Path      string         `json:"path"`
	Actor     string         `json:"actor"`
	Event     AuditEventKind `json:"event"`
	Version   int            `json:"version,omitempty"`
	Note      string         `json:"note,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// FullPath returns the canonical mount+path string.
func (e AuditTrailEntry) FullPath() string {
	return fmt.Sprintf("%s/%s", e.Mount, e.Path)
}

// Validate returns an error if the entry is missing required fields.
func (e AuditTrailEntry) Validate() error {
	if e.Mount == "" {
		return errors.New("audit trail entry: mount is required")
	}
	if e.Path == "" {
		return errors.New("audit trail entry: path is required")
	}
	if e.Actor == "" {
		return errors.New("audit trail entry: actor is required")
	}
	if !IsValidAuditEvent(e.Event) {
		return fmt.Errorf("audit trail entry: unknown event kind %q", e.Event)
	}
	return nil
}
