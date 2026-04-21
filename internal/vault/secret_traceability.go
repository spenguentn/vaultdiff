package vault

import (
	"errors"
	"fmt"
	"time"
)

// TraceabilitySource represents the origin tracking source for a secret.
type TraceabilitySource string

const (
	TraceSourceManual    TraceabilitySource = "manual"
	TraceSourcePipeline  TraceabilitySource = "pipeline"
	TraceSourceImport    TraceabilitySource = "import"
	TraceSourceGenerated TraceabilitySource = "generated"
	TraceSourceMigrated  TraceabilitySource = "migrated"
)

// IsValidTraceabilitySource returns true if s is a known traceability source.
func IsValidTraceabilitySource(s TraceabilitySource) bool {
	switch s {
	case TraceSourceManual, TraceSourcePipeline, TraceSourceImport,
		TraceSourceGenerated, TraceSourceMigrated:
		return true
	}
	return false
}

// SecretTraceability records the traceability metadata for a secret version.
type SecretTraceability struct {
	Mount      string             `json:"mount"`
	Path       string             `json:"path"`
	Version    int                `json:"version"`
	Source     TraceabilitySource `json:"source"`
	TracedBy   string             `json:"traced_by"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	TracedAt   time.Time          `json:"traced_at"`
}

// FullPath returns the canonical mount+path string.
func (t *SecretTraceability) FullPath() string {
	return fmt.Sprintf("%s/%s", t.Mount, t.Path)
}

// Validate returns an error if the traceability record is incomplete.
func (t *SecretTraceability) Validate() error {
	if t.Mount == "" {
		return errors.New("traceability: mount is required")
	}
	if t.Path == "" {
		return errors.New("traceability: path is required")
	}
	if t.Version < 1 {
		return errors.New("traceability: version must be >= 1")
	}
	if !IsValidTraceabilitySource(t.Source) {
		return fmt.Errorf("traceability: unknown source %q", t.Source)
	}
	if t.TracedBy == "" {
		return errors.New("traceability: traced_by is required")
	}
	return nil
}
