package vault

import (
	"fmt"
	"strings"
	"time"
)

// ConfidenceLevel represents the degree of confidence in a secret's validity or accuracy.
type ConfidenceLevel string

const (
	ConfidenceHigh   ConfidenceLevel = "high"
	ConfidenceMedium ConfidenceLevel = "medium"
	ConfidenceLow    ConfidenceLevel = "low"
	ConfidenceUnknown ConfidenceLevel = "unknown"
)

// IsValidConfidenceLevel returns true if the given level is a known confidence level.
func IsValidConfidenceLevel(level ConfidenceLevel) bool {
	switch level {
	case ConfidenceHigh, ConfidenceMedium, ConfidenceLow, ConfidenceUnknown:
		return true
	}
	return false
}

// SecretConfidence records the confidence level assigned to a secret at a given path.
type SecretConfidence struct {
	Mount      string          `json:"mount"`
	Path       string          `json:"path"`
	Level      ConfidenceLevel `json:"level"`
	Reason     string          `json:"reason,omitempty"`
	AssignedBy string          `json:"assigned_by"`
	AssignedAt time.Time       `json:"assigned_at"`
}

// FullPath returns the combined mount and path for the secret.
func (sc SecretConfidence) FullPath() string {
	return strings.Trim(sc.Mount, "/") + "/" + strings.Trim(sc.Path, "/")
}

// Validate checks that the SecretConfidence record has all required fields.
func (sc SecretConfidence) Validate() error {
	if sc.Mount == "" {
		return fmt.Errorf("confidence: mount is required")
	}
	if sc.Path == "" {
		return fmt.Errorf("confidence: path is required")
	}
	if !IsValidConfidenceLevel(sc.Level) {
		return fmt.Errorf("confidence: unknown level %q", sc.Level)
	}
	if sc.AssignedBy == "" {
		return fmt.Errorf("confidence: assigned_by is required")
	}
	return nil
}

// confidenceKey builds the registry key for a confidence record.
func confidenceKey(mount, path string) string {
	return strings.Trim(mount, "/") + "/" + strings.Trim(path, "/")
}

// SecretConfidenceRegistry stores confidence levels for secrets in memory.
type SecretConfidenceRegistry struct {
	entries map[string]SecretConfidence
}

// NewSecretConfidenceRegistry creates an empty SecretConfidenceRegistry.
func NewSecretConfidenceRegistry() *SecretConfidenceRegistry {
	return &SecretConfidenceRegistry{entries: make(map[string]SecretConfidence)}
}

// Set stores a confidence record, setting AssignedAt if not already set.
func (r *SecretConfidenceRegistry) Set(sc SecretConfidence) error {
	if err := sc.Validate(); err != nil {
		return err
	}
	if sc.AssignedAt.IsZero() {
		sc.AssignedAt = time.Now().UTC()
	}
	r.entries[confidenceKey(sc.Mount, sc.Path)] = sc
	return nil
}

// Get retrieves the confidence record for the given mount and path.
func (r *SecretConfidenceRegistry) Get(mount, path string) (SecretConfidence, bool) {
	v, ok := r.entries[confidenceKey(mount, path)]
	return v, ok
}

// Remove deletes the confidence record for the given mount and path.
func (r *SecretConfidenceRegistry) Remove(mount, path string) {
	delete(r.entries, confidenceKey(mount, path))
}

// All returns a copy of all stored confidence records.
func (r *SecretConfidenceRegistry) All() []SecretConfidence {
	out := make([]SecretConfidence, 0, len(r.entries))
	for _, v := range r.entries {
		out = append(out, v)
	}
	return out
}
