package vault

import (
	"testing"
	"time"
)

func TestNewSecretDelegationRegistry_NotNil(t *testing.T) {
	r := NewSecretDelegationRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestDelegationRegistry_Delegate_And_Get(t *testing.T) {
	r := NewSecretDelegationRegistry()
	d := SecretDelegation{
		Mount:      "secret",
		Path:       "app/db",
		DelegateTo: "sre-team",
		Scope:      DelegationScopeRead,
		GrantedBy:  "alice",
	}
	if err := r.Delegate(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/db", "sre-team")
	if !ok {
		t.Fatal("expected delegation to be found")
	}
	if got.GrantedBy != "alice" {
		t.Errorf("expected GrantedBy alice, got %s", got.GrantedBy)
	}
}

func TestDelegationRegistry_Delegate_SetsDelegatedAt(t *testing.T) {
	r := NewSecretDelegationRegistry()
	before := time.Now().UTC()
	d := SecretDelegation{
		Mount:      "secret",
		Path:       "app/db",
		DelegateTo: "ops",
		Scope:      DelegationScopeRead,
		GrantedBy:  "bob",
	}
	_ = r.Delegate(d)
	got, _ := r.Get("secret", "app/db", "ops")
	if got.DelegatedAt.Before(before) {
		t.Error("expected DelegatedAt to be set")
	}
}

func TestDelegationRegistry_Delegate_Invalid(t *testing.T) {
	r := NewSecretDelegationRegistry()
	if err := r.Delegate(SecretDelegation{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestDelegationRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretDelegationRegistry()
	_, ok := r.Get("secret", "missing", "nobody")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestDelegationRegistry_Revoke(t *testing.T) {
	r := NewSecretDelegationRegistry()
	d := SecretDelegation{
		Mount:      "secret",
		Path:       "app/key",
		DelegateTo: "dev",
		Scope:      DelegationScopeRead,
		GrantedBy:  "carol",
	}
	_ = r.Delegate(d)
	r.Revoke("secret", "app/key", "dev")
	_, ok := r.Get("secret", "app/key", "dev")
	if ok {
		t.Fatal("expected delegation to be revoked")
	}
}

func TestDelegationRegistry_Active_FiltersExpired(t *testing.T) {
	r := NewSecretDelegationRegistry()
	active := SecretDelegation{
		Mount:      "secret",
		Path:       "a",
		DelegateTo: "x",
		Scope:      DelegationScopeRead,
		GrantedBy:  "dave",
		ExpiresAt:  time.Now().Add(time.Hour),
	}
	expired := SecretDelegation{
		Mount:      "secret",
		Path:       "b",
		DelegateTo: "y",
		Scope:      DelegationScopeRead,
		GrantedBy:  "dave",
		ExpiresAt:  time.Now().Add(-time.Hour),
	}
	_ = r.Delegate(active)
	_ = r.Delegate(expired)
	got := r.Active()
	if len(got) != 1 {
		t.Errorf("expected 1 active delegation, got %d", len(got))
	}
}
