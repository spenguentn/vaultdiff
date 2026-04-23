package vault_test

import (
	"testing"
	"time"
)

func TestNewSecretCoverageRegistry_NotNil(t *testing.T) {
	reg := NewSecretCoverageRegistry()
	if reg == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestCoverageRegistry_Set_And_Get(t *testing.T) {
	reg := NewSecretCoverageRegistry()
	c := SecretCoverage{
		Mount:  "secret",
		Path:   "svc/api",
		Status: "full",
	}
	if err := reg.Set(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("secret", "svc/api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Status != "full" {
		t.Errorf("expected full, got %s", got.Status)
	}
}

func TestCoverageRegistry_Set_SetsRecordedAt(t *testing.T) {
	reg := NewSecretCoverageRegistry()
	before := time.Now()
	c := SecretCoverage{
		Mount:  "secret",
		Path:   "svc/token",
		Status: "partial",
	}
	_ = reg.Set(c)
	got, _ := reg.Get("secret", "svc/token")
	if got.RecordedAt.Before(before) {
		t.Error("expected RecordedAt to be set on Set")
	}
}

func TestCoverageRegistry_Set_Invalid(t *testing.T) {
	reg := NewSecretCoverageRegistry()
	c := SecretCoverage{Path: "svc/api", Status: "full"}
	if err := reg.Set(c); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestCoverageRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretCoverageRegistry()
	_, ok := reg.Get("secret", "nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestCoverageRegistry_Remove(t *testing.T) {
	reg := NewSecretCoverageRegistry()
	c := SecretCoverage{
		Mount:  "secret",
		Path:   "svc/old",
		Status: "none",
	}
	_ = reg.Set(c)
	reg.Remove("secret", "svc/old")
	_, ok := reg.Get("secret", "svc/old")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}
