package vault

import (
	"testing"
	"time"
)

func TestNewSecretTenureRegistry_NotNil(t *testing.T) {
	if NewSecretTenureRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestTenureRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretTenureRegistry()
	created := time.Now().Add(-200 * 24 * time.Hour)
	if err := r.Set("secret", "app/key", created); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ten, ok := r.Get("secret", "app/key")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if ten.Status != TenureMature {
		t.Errorf("expected TenureMature, got %s", ten.Status)
	}
}

func TestTenureRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretTenureRegistry()
	if err := r.Set("", "app/key", time.Now()); err == nil {
		t.Error("expected error for empty mount")
	}
}

func TestTenureRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretTenureRegistry()
	_, ok := r.Get("secret", "missing/key")
	if ok {
		t.Error("expected not found")
	}
}

func TestTenureRegistry_Remove(t *testing.T) {
	r := NewSecretTenureRegistry()
	_ = r.Set("secret", "app/key", time.Now().Add(-10*24*time.Hour))
	r.Remove("secret", "app/key")
	_, ok := r.Get("secret", "app/key")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestTenureRegistry_All(t *testing.T) {
	r := NewSecretTenureRegistry()
	_ = r.Set("secret", "app/key1", time.Now().Add(-10*24*time.Hour))
	_ = r.Set("secret", "app/key2", time.Now().Add(-400*24*time.Hour))
	if len(r.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.All()))
	}
}
