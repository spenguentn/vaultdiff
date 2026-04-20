package vault

import (
	"testing"
	"time"
)

func TestNewSecretConsentRegistry_NotNil(t *testing.T) {
	if NewSecretConsentRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestConsentRegistry_Grant_And_Get(t *testing.T) {
	r := NewSecretConsentRegistry()
	c := baseConsent()
	if err := r.Grant(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get(c.Mount, c.Path, c.GrantedTo)
	if !ok {
		t.Fatal("expected record to exist")
	}
	if got.GrantedTo != c.GrantedTo {
		t.Errorf("got %s, want %s", got.GrantedTo, c.GrantedTo)
	}
}

func TestConsentRegistry_Grant_SetsGrantedAt(t *testing.T) {
	r := NewSecretConsentRegistry()
	c := baseConsent()
	c.GrantedAt = time.Time{}
	_ = r.Grant(c)
	got, _ := r.Get(c.Mount, c.Path, c.GrantedTo)
	if got.GrantedAt.IsZero() {
		t.Error("expected GrantedAt to be set automatically")
	}
}

func TestConsentRegistry_Grant_Invalid(t *testing.T) {
	r := NewSecretConsentRegistry()
	c := baseConsent()
	c.Mount = ""
	if err := r.Grant(c); err == nil {
		t.Error("expected error for invalid consent")
	}
}

func TestConsentRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretConsentRegistry()
	_, ok := r.Get("secret", "missing", "alice")
	if ok {
		t.Error("expected not found")
	}
}

func TestConsentRegistry_Revoke(t *testing.T) {
	r := NewSecretConsentRegistry()
	c := baseConsent()
	_ = r.Grant(c)
	if err := r.Revoke(c.Mount, c.Path, c.GrantedTo); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := r.Get(c.Mount, c.Path, c.GrantedTo)
	if got.Status != ConsentRevoked {
		t.Errorf("expected revoked, got %s", got.Status)
	}
}

func TestConsentRegistry_Revoke_NotFound(t *testing.T) {
	r := NewSecretConsentRegistry()
	if err := r.Revoke("secret", "missing", "alice"); err == nil {
		t.Error("expected error for missing record")
	}
}

func TestConsentRegistry_Remove(t *testing.T) {
	r := NewSecretConsentRegistry()
	c := baseConsent()
	_ = r.Grant(c)
	r.Remove(c.Mount, c.Path, c.GrantedTo)
	_, ok := r.Get(c.Mount, c.Path, c.GrantedTo)
	if ok {
		t.Error("expected record to be removed")
	}
}

func TestConsentRegistry_All(t *testing.T) {
	r := NewSecretConsentRegistry()
	c1 := baseConsent()
	c2 := baseConsent()
	c2.GrantedTo = "bob"
	_ = r.Grant(c1)
	_ = r.Grant(c2)
	if got := len(r.All()); got != 2 {
		t.Errorf("expected 2 records, got %d", got)
	}
}
