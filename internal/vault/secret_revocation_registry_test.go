package vault

import (
	"testing"
	"time"
)

func baseRevocation() RevocationRecord {
	return RevocationRecord{
		Mount:     "secret",
		Path:      "app/db-password",
		RevokedBy: "ops-team",
		Reason:    "compromised",
		RevokedAt: time.Now().UTC(),
	}
}

func TestNewSecretRevocationRegistry_NotNil(t *testing.T) {
	if NewSecretRevocationRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestRevocationRegistry_Revoke_And_Get(t *testing.T) {
	reg := NewSecretRevocationRegistry()
	rec := baseRevocation()
	if err := reg.Revoke(rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get(rec.Mount, rec.Path)
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.RevokedBy != rec.RevokedBy {
		t.Errorf("revokedBy mismatch: got %s", got.RevokedBy)
	}
}

func TestRevocationRegistry_Revoke_SetsRevokedAt(t *testing.T) {
	reg := NewSecretRevocationRegistry()
	rec := baseRevocation()
	rec.RevokedAt = time.Time{}
	if err := reg.Revoke(rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := reg.Get(rec.Mount, rec.Path)
	if got.RevokedAt.IsZero() {
		t.Error("expected RevokedAt to be set automatically")
	}
}

func TestRevocationRegistry_Revoke_Invalid(t *testing.T) {
	reg := NewSecretRevocationRegistry()
	rec := baseRevocation()
	rec.Mount = ""
	if err := reg.Revoke(rec); err == nil {
		t.Error("expected error for invalid record")
	}
}

func TestRevocationRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretRevocationRegistry()
	_, ok := reg.Get("secret", "nonexistent")
	if ok {
		t.Error("expected not found")
	}
}

func TestRevocationRegistry_Remove(t *testing.T) {
	reg := NewSecretRevocationRegistry()
	rec := baseRevocation()
	_ = reg.Revoke(rec)
	reg.Remove(rec.Mount, rec.Path)
	_, ok := reg.Get(rec.Mount, rec.Path)
	if ok {
		t.Error("expected record to be removed")
	}
}

func TestRevocationRegistry_All(t *testing.T) {
	reg := NewSecretRevocationRegistry()
	r1 := baseRevocation()
	r2 := baseRevocation()
	r2.Path = "app/api-key"
	_ = reg.Revoke(r1)
	_ = reg.Revoke(r2)
	if len(reg.All()) != 2 {
		t.Errorf("expected 2 records, got %d", len(reg.All()))
	}
}
