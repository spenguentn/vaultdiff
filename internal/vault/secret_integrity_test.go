package vault

import (
	"testing"
	"time"
)

func TestIsValidIntegrityStatus_Known(t *testing.T) {
	for _, s := range []IntegrityStatus{IntegrityStatusOK, IntegrityStatusTampered, IntegrityStatusUnknown} {
		if !IsValidIntegrityStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidIntegrityStatus_Unknown(t *testing.T) {
	if IsValidIntegrityStatus("bogus") {
		t.Error("expected 'bogus' to be invalid")
	}
}

func TestSecretIntegrityRecord_FullPath(t *testing.T) {
	rec := &SecretIntegrityRecord{Mount: "/secret", Path: "/myapp/db"}
	if got := rec.FullPath(); got != "secret/myapp/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretIntegrityRecord_Validate_Valid(t *testing.T) {
	rec := &SecretIntegrityRecord{
		Mount: "secret", Path: "app/db",
		Fingerprint: "abc123", Status: IntegrityStatusOK,
		CheckedBy: "ci-bot",
	}
	if err := rec.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretIntegrityRecord_Validate_MissingMount(t *testing.T) {
	rec := &SecretIntegrityRecord{
		Path: "app/db", Fingerprint: "abc",
		Status: IntegrityStatusOK, CheckedBy: "ci",
	}
	if err := rec.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretIntegrityRecord_Validate_InvalidStatus(t *testing.T) {
	rec := &SecretIntegrityRecord{
		Mount: "secret", Path: "app", Fingerprint: "abc",
		Status: "nope", CheckedBy: "ci",
	}
	if err := rec.Validate(); err == nil {
		t.Error("expected error for invalid status")
	}
}

func TestComputeIntegrityFingerprint_Deterministic(t *testing.T) {
	data := map[string]string{"b": "2", "a": "1"}
	f1, err := ComputeIntegrityFingerprint(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	f2, _ := ComputeIntegrityFingerprint(data)
	if f1 != f2 {
		t.Error("fingerprint not deterministic")
	}
}

func TestComputeIntegrityFingerprint_NilData(t *testing.T) {
	if _, err := ComputeIntegrityFingerprint(nil); err == nil {
		t.Error("expected error for nil data")
	}
}

func TestNewSecretIntegrityRegistry_NotNil(t *testing.T) {
	if NewSecretIntegrityRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestIntegrityRegistry_Record_And_Get(t *testing.T) {
	reg := NewSecretIntegrityRegistry()
	rec := &SecretIntegrityRecord{
		Mount: "secret", Path: "app/key",
		Fingerprint: "deadbeef", Status: IntegrityStatusOK,
		CheckedBy: "audit",
	}
	if err := reg.Record(rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("secret", "app/key")
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.Fingerprint != "deadbeef" {
		t.Errorf("unexpected fingerprint: %s", got.Fingerprint)
	}
}

func TestIntegrityRegistry_Record_SetsCheckedAt(t *testing.T) {
	reg := NewSecretIntegrityRegistry()
	before := time.Now()
	rec := &SecretIntegrityRecord{
		Mount: "secret", Path: "x",
		Fingerprint: "fp", Status: IntegrityStatusOK,
		CheckedBy: "bot",
	}
	_ = reg.Record(rec)
	if rec.CheckedAt.Before(before) {
		t.Error("expected CheckedAt to be stamped")
	}
}

func TestIntegrityRegistry_Remove(t *testing.T) {
	reg := NewSecretIntegrityRegistry()
	rec := &SecretIntegrityRecord{
		Mount: "secret", Path: "rm",
		Fingerprint: "fp", Status: IntegrityStatusOK,
		CheckedBy: "bot",
	}
	_ = reg.Record(rec)
	reg.Remove("secret", "rm")
	if _, ok := reg.Get("secret", "rm"); ok {
		t.Error("expected record to be removed")
	}
}
