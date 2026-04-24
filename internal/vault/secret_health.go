package vault

import (
	"fmt"
	"time"
)

// SecretHealthStatus represents the overall health of a secret.
type SecretHealthStatus string

const (
	SecretHealthOK       SecretHealthStatus = "ok"
	SecretHealthDegraded SecretHealthStatus = "degraded"
	SecretHealthCritical SecretHealthStatus = "critical"
	SecretHealthUnknown  SecretHealthStatus = "unknown"
)

// IsValidSecretHealthStatus returns true if s is a known health status.
func IsValidSecretHealthStatus(s SecretHealthStatus) bool {
	switch s {
	case SecretHealthOK, SecretHealthDegraded, SecretHealthCritical, SecretHealthUnknown:
		return true
	}
	return false
}

// SecretHealth records the computed health of a secret at a point in time.
type SecretHealth struct {
	Mount      string             `json:"mount"`
	Path       string             `json:"path"`
	Status     SecretHealthStatus `json:"status"`
	Reason     string             `json:"reason,omitempty"`
	CheckedAt  time.Time          `json:"checked_at"`
	CheckedBy  string             `json:"checked_by"`
}

// FullPath returns the canonical mount+path identifier.
func (h SecretHealth) FullPath() string {
	return fmt.Sprintf("%s/%s", h.Mount, h.Path)
}

// IsHealthy returns true when the status is OK.
func (h SecretHealth) IsHealthy() bool {
	return h.Status == SecretHealthOK
}

// Validate returns an error if the record is incomplete or invalid.
func (h SecretHealth) Validate() error {
	if h.Mount == "" {
		return fmt.Errorf("secret health: mount is required")
	}
	if h.Path == "" {
		return fmt.Errorf("secret health: path is required")
	}
	if !IsValidSecretHealthStatus(h.Status) {
		return fmt.Errorf("secret health: unknown status %q", h.Status)
	}
	if h.CheckedBy == "" {
		return fmt.Errorf("secret health: checked_by is required")
	}
	if h.CheckedAt.IsZero() {
		return fmt.Errorf("secret health: checked_at is required")
	}
	return nil
}

// healthKey builds the registry lookup key.
func healthKey(mount, path string) string {
	return mount + "/" + path
}

// SecretHealthRegistry stores health records keyed by mount+path.
type SecretHealthRegistry struct {
	records map[string]SecretHealth
}

// NewSecretHealthRegistry returns an empty registry.
func NewSecretHealthRegistry() *SecretHealthRegistry {
	return &SecretHealthRegistry{records: make(map[string]SecretHealth)}
}

// Set validates and stores a health record, stamping CheckedAt if zero.
func (r *SecretHealthRegistry) Set(h SecretHealth) error {
	if h.CheckedAt.IsZero() {
		h.CheckedAt = time.Now().UTC()
	}
	if err := h.Validate(); err != nil {
		return err
	}
	r.records[healthKey(h.Mount, h.Path)] = h
	return nil
}

// Get retrieves a health record by mount and path.
func (r *SecretHealthRegistry) Get(mount, path string) (SecretHealth, bool) {
	v, ok := r.records[healthKey(mount, path)]
	return v, ok
}

// Remove deletes a health record.
func (r *SecretHealthRegistry) Remove(mount, path string) {
	delete(r.records, healthKey(mount, path))
}
