package vault

import "testing"

func TestNewSecretSeverityRegistry_NotNil(t *testing.T) {
	if NewSecretSeverityRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestSeverityRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretSeverityRegistry()
	s := SecretSeverity{Mount: "secret", Path: "svc/key", Level: SeverityMedium, AssignedBy: "bob"}
	if err := r.Set(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "svc/key")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.Level != SeverityMedium {
		t.Errorf("expected medium, got %s", got.Level)
	}
}

func TestSeverityRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretSeverityRegistry()
	s := SecretSeverity{Mount: "", Path: "svc/key", Level: SeverityLow, AssignedBy: "bob"}
	if err := r.Set(s); err == nil {
		t.Error("expected validation error")
	}
}

func TestSeverityRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretSeverityRegistry()
	_, ok := r.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestSeverityRegistry_Remove(t *testing.T) {
	r := NewSecretSeverityRegistry()
	s := SecretSeverity{Mount: "secret", Path: "svc/key", Level: SeverityHigh, AssignedBy: "carol"}
	_ = r.Set(s)
	r.Remove("secret", "svc/key")
	_, ok := r.Get("secret", "svc/key")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestSeverityRegistry_All(t *testing.T) {
	r := NewSecretSeverityRegistry()
	_ = r.Set(SecretSeverity{Mount: "m", Path: "p1", Level: SeverityLow, AssignedBy: "x"})
	_ = r.Set(SecretSeverity{Mount: "m", Path: "p2", Level: SeverityCritical, AssignedBy: "x"})
	if len(r.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.All()))
	}
}
