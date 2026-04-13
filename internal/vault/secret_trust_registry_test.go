package vault

import (
	"testing"
)

func TestNewSecretTrustRegistry_NotNil(t *testing.T) {
	r := NewSecretTrustRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestTrustRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretTrustRegistry()
	tr := &SecretTrust{
		Mount:      "secret",
		Path:       "app/key",
		Level:      TrustLevelVerified,
		AssignedBy: "ops",
	}
	if err := r.Set(tr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/key")
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.Level != TrustLevelVerified {
		t.Errorf("expected verified, got %s", got.Level)
	}
}

func TestTrustRegistry_Set_SetsAssignedAt(t *testing.T) {
	r := NewSecretTrustRegistry()
	tr := &SecretTrust{Mount: "secret", Path: "x", Level: TrustLevelLow, AssignedBy: "bob"}
	_ = r.Set(tr)
	got, _ := r.Get("secret", "x")
	if got.AssignedAt.IsZero() {
		t.Error("expected AssignedAt to be set")
	}
}

func TestTrustRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretTrustRegistry()
	tr := &SecretTrust{Mount: "", Path: "x", Level: TrustLevelLow, AssignedBy: "bob"}
	if err := r.Set(tr); err == nil {
		t.Error("expected validation error")
	}
}

func TestTrustRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretTrustRegistry()
	_, ok := r.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestTrustRegistry_Remove(t *testing.T) {
	r := NewSecretTrustRegistry()
	tr := &SecretTrust{Mount: "secret", Path: "rm", Level: TrustLevelHigh, AssignedBy: "alice"}
	_ = r.Set(tr)
	if !r.Remove("secret", "rm") {
		t.Error("expected remove to return true")
	}
	_, ok := r.Get("secret", "rm")
	if ok {
		t.Error("expected record to be gone")
	}
}

func TestTrustRegistry_Remove_NotFound(t *testing.T) {
	r := NewSecretTrustRegistry()
	if r.Remove("secret", "ghost") {
		t.Error("expected false for missing key")
	}
}

func TestTrustRegistry_All(t *testing.T) {
	r := NewSecretTrustRegistry()
	_ = r.Set(&SecretTrust{Mount: "m", Path: "a", Level: TrustLevelLow, AssignedBy: "x"})
	_ = r.Set(&SecretTrust{Mount: "m", Path: "b", Level: TrustLevelMedium, AssignedBy: "y"})
	if len(r.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.All()))
	}
}
