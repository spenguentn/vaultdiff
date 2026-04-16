package vault

import "fmt"

// ScopeLevel defines the operational scope of a secret.
type ScopeLevel string

const (
	ScopeGlobal  ScopeLevel = "global"
	ScopeRegional ScopeLevel = "regional"
	ScopeLocal   ScopeLevel = "local"
	ScopeTeam    ScopeLevel = "team"
)

var validScopeLevels = map[ScopeLevel]bool{
	ScopeGlobal:   true,
	ScopeRegional: true,
	ScopeLocal:    true,
	ScopeTeam:     true,
}

// IsValidScopeLevel returns true if the given scope level is recognised.
func IsValidScopeLevel(s ScopeLevel) bool {
	return validScopeLevels[s]
}

// SecretScope associates a scope level with a secret path.
type SecretScope struct {
	Mount     string     `json:"mount"`
	Path      string     `json:"path"`
	Scope     ScopeLevel `json:"scope"`
	SetBy     string     `json:"set_by"`
	Namespace string     `json:"namespace,omitempty"`
}

// FullPath returns the canonical mount+path string.
func (s SecretScope) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Mount, s.Path)
}

// Validate checks that the SecretScope has all required fields.
func (s SecretScope) Validate() error {
	if s.Mount == "" {
		return fmt.Errorf("secret scope: mount is required")
	}
	if s.Path == "" {
		return fmt.Errorf("secret scope: path is required")
	}
	if !IsValidScopeLevel(s.Scope) {
		return fmt.Errorf("secret scope: invalid scope level %q", s.Scope)
	}
	if s.SetBy == "" {
		return fmt.Errorf("secret scope: set_by is required")
	}
	return nil
}
