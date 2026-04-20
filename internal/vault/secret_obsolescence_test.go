package vault

import (
	"testing"
	"time"
)

func TestIsValidObsolescenceReason_Known(t *testing.T) {
	for _, r := range []ObsolescenceReason{
		ObsolescenceReasonSuperseded,
		ObsolescenceReasonDeprecated,
		ObsolescenceReasonUnused,
		ObsolescenceReasonExpired,
	} {
		if !IsValidObsolescenceReason(r) {
			t.Errorf("expected %q to be valid", r)
		}
	}
}

func TestIsValidObsolescenceReason_Unknown(t *testing.T) {
	if IsValidObsolescenceReason("gone") {
		t.Error("expected unknown reason to be invalid")
	}
}

func TestSecretObsolescence_FullPath(t *testing.T) {
	o := &SecretObsolescence{Mount: "secret", Path: "app/db"}
	if got := o.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretObsolescence_IsDue_False(t *testing.T) {
	future := time.Now().Add(time.Hour)
	o := &SecretObsolescence{ScheduledAt: &future}
	if o.IsDue() {
		t.Error("expected IsDue to be false for future time")
	}
}

func TestSecretObsolescence_IsDue_True(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	o := &SecretObsolescence{ScheduledAt: &past}
	if !o.IsDue() {
		t.Error("expected IsDue to be true for past time")
	}
}

func TestSecretObsolescence_IsDue_NoSchedule(t *testing.T) {
	o := &SecretObsolescence{}
	if o.IsDue() {
		t.Error("expected IsDue to be false when ScheduledAt is nil")
	}
}

func TestSecretObsolescence_Validate_Valid(t *testing.T) {
	o := &SecretObsolescence{
		Mount:    "secret",
		Path:     "app/key",
		Reason:   ObsolescenceReasonDeprecated,
		MarkedBy: "ops-team",
	}
	if err := o.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretObsolescence_Validate_MissingMount(t *testing.T) {
	o := &SecretObsolescence{Path: "app/key", Reason: ObsolescenceReasonUnused, MarkedBy: "alice"}
	if err := o.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretObsolescence_Validate_InvalidReason(t *testing.T) {
	o := &SecretObsolescence{Mount: "secret", Path: "app/key", Reason: "gone", MarkedBy: "alice"}
	if err := o.Validate(); err == nil {
		t.Error("expected error for invalid reason")
	}
}

func TestObsolescenceRegistry_Mark_And_Get(t *testing.T) {
	reg := NewSecretObsolescenceRegistry()
	o := &SecretObsolescence{
		Mount:    "kv",
		Path:     "svc/token",
		Reason:   ObsolescenceReasonSuperseded,
		MarkedBy: "bot",
	}
	if err := reg.Mark(o); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("kv", "svc/token")
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.MarkedAt.IsZero() {
		t.Error("expected MarkedAt to be set")
	}
}

func TestObsolescenceRegistry_Mark_Invalid(t *testing.T) {
	reg := NewSecretObsolescenceRegistry()
	if err := reg.Mark(&SecretObsolescence{}); err == nil {
		t.Error("expected validation error")
	}
}

func TestObsolescenceRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretObsolescenceRegistry()
	if _, ok := reg.Get("kv", "missing"); ok {
		t.Error("expected not found")
	}
}

func TestObsolescenceRegistry_Remove(t *testing.T) {
	reg := NewSecretObsolescenceRegistry()
	_ = reg.Mark(&SecretObsolescence{Mount: "kv", Path: "x", Reason: ObsolescenceReasonExpired, MarkedBy: "ci"})
	reg.Remove("kv", "x")
	if _, ok := reg.Get("kv", "x"); ok {
		t.Error("expected record to be removed")
	}
}

func TestObsolescenceRegistry_All(t *testing.T) {
	reg := NewSecretObsolescenceRegistry()
	for _, p := range []string{"a", "b", "c"} {
		_ = reg.Mark(&SecretObsolescence{Mount: "kv", Path: p, Reason: ObsolescenceReasonUnused, MarkedBy: "ci"})
	}
	if len(reg.All()) != 3 {
		t.Errorf("expected 3 entries, got %d", len(reg.All()))
	}
}
