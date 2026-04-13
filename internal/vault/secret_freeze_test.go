package vault

import (
	"testing"
	"time"
)

func TestSecretFreeze_FullPath(t *testing.T) {
	f := SecretFreeze{Mount: "secret", Path: "app/db"}
	if got := f.FullPath(); got != "secret/app/db" {
		t.Fatalf("expected secret/app/db, got %s", got)
	}
}

func TestSecretFreeze_IsExpired_False(t *testing.T) {
	future := time.Now().Add(time.Hour)
	f := SecretFreeze{ExpiresAt: &future}
	if f.IsExpired() {
		t.Fatal("expected not expired")
	}
}

func TestSecretFreeze_IsExpired_True(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	f := SecretFreeze{ExpiresAt: &past}
	if !f.IsExpired() {
		t.Fatal("expected expired")
	}
}

func TestSecretFreeze_IsExpired_NoExpiry(t *testing.T) {
	f := SecretFreeze{}
	if f.IsExpired() {
		t.Fatal("expected not expired when no expiry set")
	}
}

func TestSecretFreeze_Validate_Valid(t *testing.T) {
	f := SecretFreeze{Mount: "secret", Path: "app/key", FrozenBy: "alice", Reason: FreezeReasonManual}
	if err := f.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretFreeze_Validate_MissingMount(t *testing.T) {
	f := SecretFreeze{Path: "app/key", FrozenBy: "alice", Reason: FreezeReasonManual}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestSecretFreeze_Validate_MissingFrozenBy(t *testing.T) {
	f := SecretFreeze{Mount: "secret", Path: "app/key", Reason: FreezeReasonManual}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for missing frozen_by")
	}
}

func TestSecretFreeze_Validate_UnknownReason(t *testing.T) {
	f := SecretFreeze{Mount: "secret", Path: "app/key", FrozenBy: "alice", Reason: "unknown"}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for unknown reason")
	}
}

func TestIsValidFreezeReason_Known(t *testing.T) {
	for _, r := range []FreezeReason{FreezeReasonManual, FreezeReasonCompliance, FreezeReasonIncident} {
		if !IsValidFreezeReason(r) {
			t.Fatalf("expected %q to be valid", r)
		}
	}
}

func TestIsValidFreezeReason_Unknown(t *testing.T) {
	if IsValidFreezeReason("bogus") {
		t.Fatal("expected bogus to be invalid")
	}
}

func TestNewSecretFreezeRegistry_NotNil(t *testing.T) {
	if NewSecretFreezeRegistry() == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestFreezeRegistry_Freeze_And_Get(t *testing.T) {
	reg := NewSecretFreezeRegistry()
	f := SecretFreeze{Mount: "secret", Path: "db/pass", FrozenBy: "ops", Reason: FreezeReasonIncident}
	if err := reg.Freeze(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("secret", "db/pass")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if got.FrozenBy != "ops" {
		t.Fatalf("expected ops, got %s", got.FrozenBy)
	}
}

func TestFreezeRegistry_Freeze_SetsFrozenAt(t *testing.T) {
	reg := NewSecretFreezeRegistry()
	f := SecretFreeze{Mount: "secret", Path: "db/pass", FrozenBy: "ops", Reason: FreezeReasonManual}
	_ = reg.Freeze(f)
	got, _ := reg.Get("secret", "db/pass")
	if got.FrozenAt.IsZero() {
		t.Fatal("expected FrozenAt to be set")
	}
}

func TestFreezeRegistry_IsFrozen_True(t *testing.T) {
	reg := NewSecretFreezeRegistry()
	f := SecretFreeze{Mount: "kv", Path: "svc/token", FrozenBy: "ci", Reason: FreezeReasonCompliance}
	_ = reg.Freeze(f)
	if !reg.IsFrozen("kv", "svc/token") {
		t.Fatal("expected path to be frozen")
	}
}

func TestFreezeRegistry_IsFrozen_NotFound(t *testing.T) {
	reg := NewSecretFreezeRegistry()
	if reg.IsFrozen("kv", "missing/path") {
		t.Fatal("expected path not to be frozen")
	}
}

func TestFreezeRegistry_Unfreeze(t *testing.T) {
	reg := NewSecretFreezeRegistry()
	f := SecretFreeze{Mount: "kv", Path: "svc/key", FrozenBy: "admin", Reason: FreezeReasonManual}
	_ = reg.Freeze(f)
	if !reg.Unfreeze("kv", "svc/key") {
		t.Fatal("expected unfreeze to return true")
	}
	if _, ok := reg.Get("kv", "svc/key"); ok {
		t.Fatal("expected record to be removed")
	}
}

func TestFreezeRegistry_Unfreeze_NotFound(t *testing.T) {
	reg := NewSecretFreezeRegistry()
	if reg.Unfreeze("kv", "ghost/path") {
		t.Fatal("expected false for missing record")
	}
}

func TestFreezeRegistry_All(t *testing.T) {
	reg := NewSecretFreezeRegistry()
	for _, p := range []string{"a", "b", "c"} {
		_ = reg.Freeze(SecretFreeze{Mount: "m", Path: p, FrozenBy: "x", Reason: FreezeReasonManual})
	}
	if len(reg.All()) != 3 {
		t.Fatalf("expected 3 records, got %d", len(reg.All()))
	}
}
