package vault

import (
	"testing"
	"time"
)

var baseAuditEntry = AuditTrailEntry{
	Mount:  "secret",
	Path:   "app/db",
	Actor:  "alice",
	Event:  AuditEventWrite,
	Version: 3,
}

func TestAuditTrailEntry_FullPath(t *testing.T) {
	if got := baseAuditEntry.FullPath(); got != "secret/app/db" {
		t.Fatalf("expected secret/app/db, got %s", got)
	}
}

func TestAuditTrailEntry_Validate_Valid(t *testing.T) {
	if err := baseAuditEntry.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuditTrailEntry_Validate_MissingMount(t *testing.T) {
	e := baseAuditEntry
	e.Mount = ""
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestAuditTrailEntry_Validate_MissingActor(t *testing.T) {
	e := baseAuditEntry
	e.Actor = ""
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for missing actor")
	}
}

func TestAuditTrailEntry_Validate_UnknownEvent(t *testing.T) {
	e := baseAuditEntry
	e.Event = "explode"
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for unknown event kind")
	}
}

func TestIsValidAuditEvent_Known(t *testing.T) {
	for _, k := range []AuditEventKind{AuditEventRead, AuditEventWrite, AuditEventDelete, AuditEventPromote, AuditEventRollback} {
		if !IsValidAuditEvent(k) {
			t.Fatalf("expected %s to be valid", k)
		}
	}
}

func TestIsValidAuditEvent_Unknown(t *testing.T) {
	if IsValidAuditEvent("unknown") {
		t.Fatal("expected unknown to be invalid")
	}
}

func TestNewSecretAuditTrailRegistry_NotNil(t *testing.T) {
	if NewSecretAuditTrailRegistry() == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestAuditTrailRegistry_Record_And_Get(t *testing.T) {
	r := NewSecretAuditTrailRegistry()
	if err := r.Record(baseAuditEntry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := r.Get("secret", "app/db")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestAuditTrailRegistry_Record_SetsTimestamp(t *testing.T) {
	r := NewSecretAuditTrailRegistry()
	e := baseAuditEntry
	e.Timestamp = time.Time{}
	_ = r.Record(e)
	got := r.Get("secret", "app/db")
	if got[0].Timestamp.IsZero() {
		t.Fatal("expected timestamp to be set automatically")
	}
}

func TestAuditTrailRegistry_Record_Invalid(t *testing.T) {
	r := NewSecretAuditTrailRegistry()
	e := baseAuditEntry
	e.Actor = ""
	if err := r.Record(e); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestAuditTrailRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretAuditTrailRegistry()
	if got := r.Get("secret", "missing"); len(got) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(got))
	}
}

func TestAuditTrailRegistry_Clear(t *testing.T) {
	r := NewSecretAuditTrailRegistry()
	_ = r.Record(baseAuditEntry)
	r.Clear("secret", "app/db")
	if got := r.Get("secret", "app/db"); len(got) != 0 {
		t.Fatal("expected entries to be cleared")
	}
}

func TestAuditTrailRegistry_Len(t *testing.T) {
	r := NewSecretAuditTrailRegistry()
	_ = r.Record(baseAuditEntry)
	e2 := baseAuditEntry
	e2.Event = AuditEventRead
	_ = r.Record(e2)
	if r.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", r.Len())
	}
}
