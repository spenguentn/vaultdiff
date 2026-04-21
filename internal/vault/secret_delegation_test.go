package vault

import (
	"testing"
	"time"
)

func TestIsValidDelegationScope_Known(t *testing.T) {
	for _, s := range []DelegationScope{DelegationScopeRead, DelegationScopeWrite, DelegationScopeAdmin} {
		if !IsValidDelegationScope(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidDelegationScope_Unknown(t *testing.T) {
	if IsValidDelegationScope("superuser") {
		t.Error("expected 'superuser' to be invalid")
	}
}

func TestSecretDelegation_FullPath(t *testing.T) {
	d := SecretDelegation{Mount: "secret", Path: "app/key"}
	if got := d.FullPath(); got != "secret/app/key" {
		t.Errorf("unexpected FullPath: %q", got)
	}
}

func TestSecretDelegation_IsExpired_False(t *testing.T) {
	d := SecretDelegation{ExpiresAt: time.Now().Add(time.Hour)}
	if d.IsExpired() {
		t.Error("expected delegation to not be expired")
	}
}

func TestSecretDelegation_IsExpired_True(t *testing.T) {
	d := SecretDelegation{ExpiresAt: time.Now().Add(-time.Minute)}
	if !d.IsExpired() {
		t.Error("expected delegation to be expired")
	}
}

func TestSecretDelegation_IsExpired_NoExpiry(t *testing.T) {
	d := SecretDelegation{}
	if d.IsExpired() {
		t.Error("expected delegation with no expiry to not be expired")
	}
}

func TestSecretDelegation_Validate_Valid(t *testing.T) {
	d := SecretDelegation{
		Mount:       "secret",
		Path:        "app/db",
		DelegatedTo: "alice",
		GrantedBy:   "admin",
		Scope:       DelegationScopeRead,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	if err := d.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretDelegation_Validate_MissingMount(t *testing.T) {
	d := SecretDelegation{Path: "x", DelegatedTo: "a", GrantedBy: "b", Scope: DelegationScopeRead}
	if err := d.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretDelegation_Validate_InvalidScope(t *testing.T) {
	d := SecretDelegation{
		Mount: "secret", Path: "x", DelegatedTo: "a", GrantedBy: "b",
		Scope: "unknown",
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for invalid scope")
	}
}

func TestSecretDelegation_Validate_PastExpiry(t *testing.T) {
	d := SecretDelegation{
		Mount: "secret", Path: "x", DelegatedTo: "a", GrantedBy: "b",
		Scope:     DelegationScopeWrite,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for past expiry")
	}
}
