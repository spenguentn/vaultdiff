package vault

import (
	"testing"
)

func TestIsValidScopeLevel_Known(t *testing.T) {
	for _, lvl := range []ScopeLevel{ScopeGlobal, ScopeRegional, ScopeLocal, ScopeTeam} {
		if !IsValidScopeLevel(lvl) {
			t.Errorf("expected %q to be valid", lvl)
		}
	}
}

func TestIsValidScopeLevel_Unknown(t *testing.T) {
	if IsValidScopeLevel("universe") {
		t.Error("expected 'universe' to be invalid")
	}
}

func TestSecretScope_FullPath(t *testing.T) {
	s := SecretScope{Mount: "secret", Path: "app/db"}
	if got := s.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretScope_Validate_Valid(t *testing.T) {
	s := SecretScope{Mount: "secret", Path: "app/db", Scope: ScopeGlobal, SetBy: "admin"}
	if err := s.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretScope_Validate_MissingMount(t *testing.T) {
	s := SecretScope{Path: "app/db", Scope: ScopeLocal, SetBy: "admin"}
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretScope_Validate_MissingPath(t *testing.T) {
	s := SecretScope{Mount: "secret", Scope: ScopeLocal, SetBy: "admin"}
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretScope_Validate_InvalidScope(t *testing.T) {
	s := SecretScope{Mount: "secret", Path: "app/db", Scope: "planetary", SetBy: "admin"}
	if err := s.Validate(); err == nil {
		t.Error("expected error for invalid scope level")
	}
}

func TestSecretScope_Validate_MissingSetBy(t *testing.T) {
	s := SecretScope{Mount: "secret", Path: "app/db", Scope: ScopeTeam}
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing set_by")
	}
}
