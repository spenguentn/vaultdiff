package vault

import (
	"testing"
	"time"
)

var baseTTLPolicy = TTLPolicy{
	Mount:     "secret",
	Path:      "app/db",
	TTL:       24 * time.Hour,
	CreatedBy: "alice",
}

func TestTTLPolicy_FullPath(t *testing.T) {
	p := baseTTLPolicy
	if got := p.FullPath(); got != "secret/app/db" {
		t.Errorf("expected secret/app/db, got %s", got)
	}
}

func TestTTLPolicy_ExpiresAt(t *testing.T) {
	p := baseTTLPolicy
	p.CreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	if !p.ExpiresAt().Equal(expected) {
		t.Errorf("unexpected ExpiresAt: %v", p.ExpiresAt())
	}
}

func TestTTLPolicy_IsExpired_False(t *testing.T) {
	p := baseTTLPolicy
	p.CreatedAt = time.Now().UTC()
	if p.IsExpired() {
		t.Error("expected policy not to be expired")
	}
}

func TestTTLPolicy_IsExpired_True(t *testing.T) {
	p := baseTTLPolicy
	p.CreatedAt = time.Now().UTC().Add(-48 * time.Hour)
	if !p.IsExpired() {
		t.Error("expected policy to be expired")
	}
}

func TestTTLPolicy_Validate_Valid(t *testing.T) {
	if err := baseTTLPolicy.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTTLPolicy_Validate_MissingMount(t *testing.T) {
	p := baseTTLPolicy
	p.Mount = ""
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestTTLPolicy_Validate_ZeroTTL(t *testing.T) {
	p := baseTTLPolicy
	p.TTL = 0
	if err := p.Validate(); err == nil {
		t.Error("expected error for zero TTL")
	}
}

func TestNewSecretTTLRegistry_NotNil(t *testing.T) {
	if NewSecretTTLRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestTTLRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretTTLRegistry()
	if err := r.Set(baseTTLPolicy); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := r.Get(baseTTLPolicy.Mount, baseTTLPolicy.Path)
	if !ok {
		t.Fatal("expected policy to be found")
	}
	if p.CreatedBy != "alice" {
		t.Errorf("unexpected CreatedBy: %s", p.CreatedBy)
	}
}

func TestTTLRegistry_Set_SetsCreatedAt(t *testing.T) {
	r := NewSecretTTLRegistry()
	_ = r.Set(baseTTLPolicy)
	p, _ := r.Get(baseTTLPolicy.Mount, baseTTLPolicy.Path)
	if p.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestTTLRegistry_Remove(t *testing.T) {
	r := NewSecretTTLRegistry()
	_ = r.Set(baseTTLPolicy)
	r.Remove(baseTTLPolicy.Mount, baseTTLPolicy.Path)
	if _, ok := r.Get(baseTTLPolicy.Mount, baseTTLPolicy.Path); ok {
		t.Error("expected policy to be removed")
	}
}

func TestTTLRegistry_Expired(t *testing.T) {
	r := NewSecretTTLRegistry()
	expired := baseTTLPolicy
	expired.CreatedAt = time.Now().UTC().Add(-48 * time.Hour)
	_ = r.Set(expired)

	fresh := baseTTLPolicy
	fresh.Path = "app/fresh"
	fresh.CreatedAt = time.Now().UTC()
	_ = r.Set(fresh)

	results := r.Expired()
	if len(results) != 1 {
		t.Errorf("expected 1 expired policy, got %d", len(results))
	}
}
