package vault

import (
	"testing"
	"time"
)

func TestNewSecretRiskRegistry_NotNil(t *testing.T) {
	r := NewSecretRiskRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestRiskRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretRiskRegistry()
	risk := SecretRisk{Mount: "secret", Path: "app/db", Level: "high"}
	if err := r.Set(risk); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/db")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.Level != "high" {
		t.Errorf("expected level 'high', got %q", got.Level)
	}
}

func TestRiskRegistry_Set_SetsAssessedAt(t *testing.T) {
	r := NewSecretRiskRegistry()
	before := time.Now().UTC()
	risk := SecretRisk{Mount: "secret", Path: "app/key", Level: "low"}
	_ = r.Set(risk)
	got, _ := r.Get("secret", "app/key")
	if got.AssessedAt.Before(before) {
		t.Errorf("AssessedAt not set correctly: %v", got.AssessedAt)
	}
}

func TestRiskRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretRiskRegistry()
	risk := SecretRisk{Mount: "", Path: "app/db", Level: "high"}
	if err := r.Set(risk); err == nil {
		t.Fatal("expected validation error for missing mount")
	}
}

func TestRiskRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretRiskRegistry()
	_, ok := r.Get("secret", "nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestRiskRegistry_Remove(t *testing.T) {
	r := NewSecretRiskRegistry()
	risk := SecretRisk{Mount: "secret", Path: "app/db", Level: "medium"}
	_ = r.Set(risk)
	r.Remove("secret", "app/db")
	_, ok := r.Get("secret", "app/db")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestRiskRegistry_All_ReturnsEntries(t *testing.T) {
	r := NewSecretRiskRegistry()
	_ = r.Set(SecretRisk{Mount: "secret", Path: "a", Level: "low"})
	_ = r.Set(SecretRisk{Mount: "secret", Path: "b", Level: "high"})
	all := r.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
