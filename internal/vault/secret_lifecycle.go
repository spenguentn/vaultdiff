package vault

import (
	"errors"
	"fmt"
	"time"
)

// LifecycleStage represents the current stage of a secret's lifecycle.
type LifecycleStage string

const (
	LifecycleStageActive     LifecycleStage = "active"
	LifecycleStageDeprecated LifecycleStage = "deprecated"
	LifecycleStageRetired    LifecycleStage = "retired"
	LifecycleStagePending    LifecycleStage = "pending"
)

// IsValidLifecycleStage returns true if the stage is a recognised value.
func IsValidLifecycleStage(s LifecycleStage) bool {
	switch s {
	case LifecycleStageActive, LifecycleStageDeprecated, LifecycleStageRetired, LifecycleStagePending:
		return true
	}
	return false
}

// SecretLifecycle records the lifecycle state of a secret at a given path.
type SecretLifecycle struct {
	Mount       string         `json:"mount"`
	Path        string         `json:"path"`
	Stage       LifecycleStage `json:"stage"`
	ManagedBy   string         `json:"managed_by"`
	TransitionAt time.Time     `json:"transition_at,omitempty"`
	Reason      string         `json:"reason,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

// FullPath returns the canonical mount+path identifier.
func (s *SecretLifecycle) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Mount, s.Path)
}

// IsPastTransition reports whether the scheduled transition time has elapsed.
func (s *SecretLifecycle) IsPastTransition() bool {
	if s.TransitionAt.IsZero() {
		return false
	}
	return time.Now().UTC().After(s.TransitionAt)
}

// Validate checks that the lifecycle record is well-formed.
func (s *SecretLifecycle) Validate() error {
	if s.Mount == "" {
		return errors.New("lifecycle: mount is required")
	}
	if s.Path == "" {
		return errors.New("lifecycle: path is required")
	}
	if !IsValidLifecycleStage(s.Stage) {
		return fmt.Errorf("lifecycle: unknown stage %q", s.Stage)
	}
	if s.ManagedBy == "" {
		return errors.New("lifecycle: managed_by is required")
	}
	return nil
}
