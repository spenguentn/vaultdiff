package vault

import (
	"testing"
	"time"
)

func TestNewSecretExpiryRegistry_NotNil(t *testing.T) {
	r := NewSecretExpiryRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestExpiryRegistry_Register_And_Get(t *testing.T) {
	r := NewSecretExpiryRegistry()
	if err := r.Register(baseExpiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get(baseExpiry.Mount, baseExpiry.Path)
	if !ok {
		t.Fatal("expected policy to be found")
	}
	if got.Owner != baseExpiry.Owner {
		t.Fatalf("got owner %q, want %q", got.Owner, baseExpiry.Owner)
	}
}

func TestExpiryRegistry_Register_Invalid(t *testing.T) {
	r := NewSecretExpiryRegistry()
	p := baseExpiry
	p.Mount = ""
	if err := r.Register(p); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestExpiryRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretExpiryRegistry()
	_, ok := r.Get("secret", "missing/path")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestExpiryRegistry_Remove(t *testing.T) {
	r := NewSecretExpiryRegistry()
	_ = r.Register(baseExpiry)
	if err := r.Remove(baseExpiry.Mount, baseExpiry.Path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Count() != 0 {
		t.Fatal("expected empty registry after remove")
	}
}

func TestExpiryRegistry_Remove_NotFound(t *testing.T) {
	r := NewSecretExpiryRegistry()
	if err := r.Remove("secret", "ghost/path"); err == nil {
		t.Fatal("expected error removing non-existent policy")
	}
}

func TestExpiryRegistry_CheckAll_Expired(t *testing.T) {
	r := NewSecretExpiryRegistry()
	p := baseExpiry
	p.ExpiresAt = time.Now().Add(-1 * time.Hour)
	_ = r.Register(p)
	statuses := r.CheckAll(time.Now())
	if len(statuses) != 1 {
		t.Fatalf("expected 1 status, got %d", len(statuses))
	}
	if !statuses[0].Expired {
		t.Fatal("expected expired status")
	}
}

func TestExpiryRegistry_Count(t *testing.T) {
	r := NewSecretExpiryRegistry()
	_ = r.Register(baseExpiry)
	if r.Count() != 1 {
		t.Fatalf("expected count 1, got %d", r.Count())
	}
}
