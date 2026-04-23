package vault

import (
	"fmt"
	"time"
)

// ScopeLevel represents the visibility boundary of a secret.
type ScopeLevel string

const (
	ScopeLevelLocal  ScopeLevel = "local"
	ScopeLevelTeam   ScopeLevel = "team"
	ScopeLevelGlobal ScopeLevel = "global"
)

// IsValidScopeLevel reports whether s is a recognised scope level.
func IsValidScopeLevel(s ScopeLevel) bool {
	switch s {
	case ScopeLevelLocal, ScopeLevelTeam, ScopeLevelGlobal:
		return true
	}
	return false
}

// SecretScope records the scope assignment for a secret.
type SecretScope struct {
	Mount      string     `json:"mount"`
	Path       string     `json:"path"`
	Level      ScopeLevel `json:"level"`
	Owner      string     `json:"owner"`
	AssignedAt time.Time  `json:"assigned_at"`
}

// FullPath returns the canonical mount/path identifier.
func (s SecretScope) FullPath() string {
	return s.Mount + "/" + s.Path
}

// Validate returns an error if the scope entry is incomplete or invalid.
func (s SecretScope) Validate() error {
	if s.Mount == "" {
		return fmt.Errorf("secret scope: mount is required")
	}
	if s.Path == "" {
		return fmt.Errorf("secret scope: path is required")
	}
	if s.Owner == "" {
		return fmt.Errorf("secret scope: owner is required")
	}
	if !IsValidScopeLevel(s.Level) {
		return fmt.Errorf("secret scope: unknown level %q", s.Level)
	}
	return nil
}
