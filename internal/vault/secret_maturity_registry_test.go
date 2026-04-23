package vault_test

import (
	"testing"
	"time"
)

func TestNewSecretMaturityRegistry_NotNil(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	if reg == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestMaturityRegistry_Set_And_Get(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	m := SecretMaturity{
		Mount: "secret",
		Path:  "app/db",
		Level: "stable",
	}
	if err := reg.Set(m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("secret", "app/db")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Level != "stable" {
		t.Errorf("expected stable, got %s", got.Level)
	}
}

func TestMaturityRegistry_Set_SetsAssessedAt(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	before := time.Now()
	m := SecretMaturity{
		Mount: "secret",
		Path:  "app/key",
		Level: "experimental",
	}
	_ = reg.Set(m)
	got, _ := reg.Get("secret", "app/key")
	if got.AssessedAt.Before(before) {
		t.Error("expected AssessedAt to be set on Set")
	}
}

func TestMaturityRegistry_Set_Invalid(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	m := SecretMaturity{Path: "app/key", Level: "stable"}
	if err := reg.Set(m); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestMaturityRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	_, ok := reg.Get("secret", "missing/path")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestMaturityRegistry_Remove(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	m := SecretMaturity{
		Mount: "secret",
		Path:  "app/remove",
		Level: "deprecated",
	}
	_ = reg.Set(m)
	reg.Remove("secret", "app/remove")
	_, ok := reg.Get("secret", "app/remove")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}
