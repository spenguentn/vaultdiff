package vault

import (
	"fmt"
	"time"
)

// AccessEventType represents the type of access event on a secret.
type AccessEventType string

const (
	AccessEventRead   AccessEventType = "read"
	AccessEventWrite  AccessEventType = "write"
	AccessEventDelete AccessEventType = "delete"
	AccessEventList   AccessEventType = "list"
)

// SecretAccessEntry records a single access event for a secret.
type SecretAccessEntry struct {
	Mount     string          `json:"mount"`
	Path      string          `json:"path"`
	EventType AccessEventType `json:"event_type"`
	Actor     string          `json:"actor"`
	Namespace string          `json:"namespace,omitempty"`
	Version   int             `json:"version,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// FullPath returns the combined mount and path.
func (e *SecretAccessEntry) FullPath() string {
	return fmt.Sprintf("%s/%s", e.Mount, e.Path)
}

// Validate checks that required fields are present.
func (e *SecretAccessEntry) Validate() error {
	if e.Mount == "" {
		return fmt.Errorf("access entry: mount is required")
	}
	if e.Path == "" {
		return fmt.Errorf("access entry: path is required")
	}
	if e.Actor == "" {
		return fmt.Errorf("access entry: actor is required")
	}
	if e.EventType == "" {
		return fmt.Errorf("access entry: event_type is required")
	}
	return nil
}

// IsValidEventType reports whether the given event type is recognised.
func IsValidEventType(t AccessEventType) bool {
	switch t {
	case AccessEventRead, AccessEventWrite, AccessEventDelete, AccessEventList:
		return true
	}
	return false
}
