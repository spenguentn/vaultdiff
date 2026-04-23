package vault

import (
	"testing"
	"time"
)

func TestNewSecretScopeRegistry_NotNil(t *testing.T) {
	r := NewSecretScopeRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestScopeRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretScopeRegistry()
	scope := SecretScope{
		Mount: "secret",
		Path:  "app/db",
		Level: ScopeLevelGlobal,
		Owner: "team-platform",
	}
	if err := r.Set(scope); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/db")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.Level != ScopeLevelGlobal {
		t.Errorf("expected level %q, got %q", ScopeLevelGlobal, got.Level)
	}
}

func TestScopeRegistry_Set_SetsAssignedAt(t *testing.T) {
	r := NewSecretScopeRegistry()
	before := time.Now()
	scope := SecretScope{
		Mount: "secret",
		Path:  "app/key",
		Level: ScopeLevelLocal,
		Owner: "team-ops",
	}
	_ = r.Set(scope)
	got, _ := r.Get("secret", "app/key")
	if got.AssignedAt.Before(before) {
		t.Error("expected AssignedAt to be set to current time")
	}
}

func TestScopeRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretScopeRegistry()
	scope := SecretScope{} // missing required fields
	if err := r.Set(scope); err == nil {
		t.Error("expected validation error")
	}
}

func TestScopeRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretScopeRegistry()
	_, ok := r.Get("secret", "nonexistent")
	if ok {
		t.Error("expected not found")
	}
}

func TestScopeRegistry_Remove(t *testing.T) {
	r := NewSecretScopeRegistry()
	scope := SecretScope{
		Mount: "secret",
		Path:  "app/token",
		Level: ScopeLevelTeam,
		Owner: "team-sec",
	}
	_ = r.Set(scope)
	r.Remove("secret", "app/token")
	_, ok := r.Get("secret", "app/token")
	if ok {
		t.Error("expected entry to be removed")
	}
}
