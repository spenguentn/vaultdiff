package vault

import (
	"testing"
	"time"
)

func TestIsValidPriorityLevel_Known(t *testing.T) {
	for _, lvl := range []PriorityLevel{PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical} {
		if !IsValidPriorityLevel(lvl) {
			t.Errorf("expected %q to be valid", lvl)
		}
	}
}

func TestIsValidPriorityLevel_Unknown(t *testing.T) {
	if IsValidPriorityLevel("urgent") {
		t.Error("expected 'urgent' to be invalid")
	}
}

func TestSecretPriority_FullPath(t *testing.T) {
	p := SecretPriority{Mount: "secret", Path: "app/db"}
	if got := p.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretPriority_Validate_Valid(t *testing.T) {
	p := SecretPriority{
		Mount:      "secret",
		Path:       "app/db",
		Level:      PriorityHigh,
		AssignedBy: "ops-team",
	}
	if err := p.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretPriority_Validate_MissingMount(t *testing.T) {
	p := SecretPriority{Path: "app/db", Level: PriorityLow, AssignedBy: "alice"}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretPriority_Validate_InvalidLevel(t *testing.T) {
	p := SecretPriority{Mount: "secret", Path: "app/db", Level: "extreme", AssignedBy: "alice"}
	if err := p.Validate(); err == nil {
		t.Error("expected error for invalid level")
	}
}

func TestNewSecretPriorityRegistry_NotNil(t *testing.T) {
	if NewSecretPriorityRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestPriorityRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretPriorityRegistry()
	p := SecretPriority{Mount: "secret", Path: "svc/key", Level: PriorityCritical, AssignedBy: "bob"}
	if err := r.Set(p); err != nil {
		t.Fatalf("unexpected Set error: %v", err)
	}
	got, err := r.Get("secret", "svc/key")
	if err != nil {
		t.Fatalf("unexpected Get error: %v", err)
	}
	if got.Level != PriorityCritical {
		t.Errorf("expected critical, got %s", got.Level)
	}
}

func TestPriorityRegistry_Set_SetsAssignedAt(t *testing.T) {
	r := NewSecretPriorityRegistry()
	p := SecretPriority{Mount: "secret", Path: "x", Level: PriorityLow, AssignedBy: "carol"}
	_ = r.Set(p)
	got, _ := r.Get("secret", "x")
	if got.AssignedAt.IsZero() {
		t.Error("expected AssignedAt to be stamped")
	}
}

func TestPriorityRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretPriorityRegistry()
	if err := r.Set(SecretPriority{}); err == nil {
		t.Error("expected validation error")
	}
}

func TestPriorityRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretPriorityRegistry()
	if _, err := r.Get("secret", "missing"); err == nil {
		t.Error("expected not-found error")
	}
}

func TestPriorityRegistry_Remove(t *testing.T) {
	r := NewSecretPriorityRegistry()
	p := SecretPriority{Mount: "m", Path: "p", Level: PriorityMedium, AssignedBy: "dave", AssignedAt: time.Now()}
	_ = r.Set(p)
	r.Remove("m", "p")
	if _, err := r.Get("m", "p"); err == nil {
		t.Error("expected entry to be removed")
	}
}

func TestPriorityRegistry_All(t *testing.T) {
	r := NewSecretPriorityRegistry()
	_ = r.Set(SecretPriority{Mount: "m", Path: "a", Level: PriorityLow, AssignedBy: "x"})
	_ = r.Set(SecretPriority{Mount: "m", Path: "b", Level: PriorityHigh, AssignedBy: "y"})
	if len(r.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.All()))
	}
}
