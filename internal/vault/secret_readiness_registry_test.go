package vault

import (
	"testing"
	"time"
)

func TestNewSecretReadinessRegistry_NotNil(t *testing.T) {
	r := NewSecretReadinessRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestReadinessRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretReadinessRegistry()
	rec := SecretReadiness{
		Mount:  "secret",
		Path:   "app/db",
		Status: ReadinessStatusReady,
		Reason: "all checks passed",
	}
	if err := r.Set(rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/db")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if got.Status != ReadinessStatusReady {
		t.Errorf("expected status %q, got %q", ReadinessStatusReady, got.Status)
	}
}

func TestReadinessRegistry_Set_SetsAssessedAt(t *testing.T) {
	r := NewSecretReadinessRegistry()
	before := time.Now().UTC()
	rec := SecretReadiness{
		Mount:  "secret",
		Path:   "app/key",
		Status: ReadinessStatusNotReady,
		Reason: "pending rotation",
	}
	if err := r.Set(rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := r.Get("secret", "app/key")
	if got.AssessedAt.Before(before) {
		t.Errorf("AssessedAt not stamped correctly: %v", got.AssessedAt)
	}
}

func TestReadinessRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretReadinessRegistry()
	rec := SecretReadiness{} // missing required fields
	if err := r.Set(rec); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestReadinessRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretReadinessRegistry()
	_, ok := r.Get("secret", "missing/path")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestReadinessRegistry_Remove(t *testing.T) {
	r := NewSecretReadinessRegistry()
	rec := SecretReadiness{
		Mount:  "secret",
		Path:   "app/remove",
		Status: ReadinessStatusReady,
		Reason: "ok",
	}
	_ = r.Set(rec)
	r.Remove("secret", "app/remove")
	_, ok := r.Get("secret", "app/remove")
	if ok {
		t.Fatal("expected record to be removed")
	}
}

func TestReadinessRegistry_All(t *testing.T) {
	r := NewSecretReadinessRegistry()
	for _, path := range []string{"a", "b", "c"} {
		_ = r.Set(SecretReadiness{
			Mount:  "secret",
			Path:   path,
			Status: ReadinessStatusReady,
			Reason: "ok",
		})
	}
	if got := len(r.All()); got != 3 {
		t.Errorf("expected 3 records, got %d", got)
	}
}
