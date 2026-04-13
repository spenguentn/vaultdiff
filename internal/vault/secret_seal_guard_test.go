package vault

import (
	"testing"
	"time"
)

func TestSealGuardEntry_FullPath(t *testing.T) {
	e := SealGuardEntry{Mount: "secret", Path: "db/password"}
	if got := e.FullPath(); got != "secret/db/password" {
		t.Errorf("expected secret/db/password, got %s", got)
	}
}

func TestSealGuardEntry_Validate_Valid(t *testing.T) {
	e := SealGuardEntry{
		Mount:     "secret",
		Path:      "api/key",
		Action:    SealGuardActionRead,
		BlockedAt: time.Now(),
	}
	if err := e.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSealGuardEntry_Validate_MissingMount(t *testing.T) {
	e := SealGuardEntry{Path: "api/key", Action: SealGuardActionWrite, BlockedAt: time.Now()}
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSealGuardEntry_Validate_MissingPath(t *testing.T) {
	e := SealGuardEntry{Mount: "secret", Action: SealGuardActionDelete, BlockedAt: time.Now()}
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSealGuardEntry_Validate_MissingAction(t *testing.T) {
	e := SealGuardEntry{Mount: "secret", Path: "api/key", BlockedAt: time.Now()}
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing action")
	}
}

func TestIsValidSealGuardAction_Known(t *testing.T) {
	for _, a := range []SealGuardAction{SealGuardActionRead, SealGuardActionWrite, SealGuardActionDelete} {
		if !IsValidSealGuardAction(a) {
			t.Errorf("expected %s to be valid", a)
		}
	}
}

func TestIsValidSealGuardAction_Unknown(t *testing.T) {
	if IsValidSealGuardAction("rotate") {
		t.Error("expected rotate to be invalid")
	}
}

func TestNewSealGuardRegistry_NotNil(t *testing.T) {
	if NewSealGuardRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestSealGuardRegistry_Record_And_Get(t *testing.T) {
	r := NewSealGuardRegistry()
	e := SealGuardEntry{
		Mount:     "kv",
		Path:      "svc/token",
		Action:    SealGuardActionRead,
		BlockedAt: time.Now(),
		Reason:    "vault sealed",
	}
	if err := r.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, ok := r.Get("kv", "svc/token")
	if !ok || len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestSealGuardRegistry_Record_SetsBlockedAt(t *testing.T) {
	r := NewSealGuardRegistry()
	e := SealGuardEntry{Mount: "kv", Path: "x", Action: SealGuardActionWrite}
	if err := r.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := r.Get("kv", "x")
	if entries[0].BlockedAt.IsZero() {
		t.Error("expected BlockedAt to be set automatically")
	}
}

func TestSealGuardRegistry_Count(t *testing.T) {
	r := NewSealGuardRegistry()
	for i := 0; i < 3; i++ {
		_ = r.Record(SealGuardEntry{
			Mount: "kv", Path: "p", Action: SealGuardActionDelete, BlockedAt: time.Now(),
		})
	}
	if r.Count() != 3 {
		t.Errorf("expected count 3, got %d", r.Count())
	}
}

func TestSealGuardRegistry_Get_NotFound(t *testing.T) {
	r := NewSealGuardRegistry()
	_, ok := r.Get("kv", "missing")
	if ok {
		t.Error("expected not found")
	}
}
