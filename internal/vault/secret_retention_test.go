package vault

import (
	"testing"
	"time"
)

var baseRetention = RetentionPolicy{
	Mount:     "secret",
	Path:      "myapp/db",
	Duration:  30,
	Unit:      RetentionUnitDays,
	CreatedBy: "admin",
	CreatedAt: time.Now().UTC(),
}

func TestRetentionPolicy_FullPath(t *testing.T) {
	p := baseRetention
	if got := p.FullPath(); got != "secret/myapp/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestRetentionPolicy_Validate_Valid(t *testing.T) {
	if err := baseRetention.Validate(); err != nil {
		t.Fatalf("expected valid, got: %v", err)
	}
}

func TestRetentionPolicy_Validate_MissingMount(t *testing.T) {
	p := baseRetention
	p.Mount = ""
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestRetentionPolicy_Validate_ZeroDuration(t *testing.T) {
	p := baseRetention
	p.Duration = 0
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestRetentionPolicy_Validate_InvalidUnit(t *testing.T) {
	p := baseRetention
	p.Unit = "years"
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for invalid unit")
	}
}

func TestRetentionPolicy_IsExpired_False(t *testing.T) {
	p := baseRetention
	p.CreatedAt = time.Now().UTC()
	if p.IsExpired() {
		t.Fatal("policy should not be expired")
	}
}

func TestRetentionPolicy_IsExpired_True(t *testing.T) {
	p := baseRetention
	p.CreatedAt = time.Now().UTC().AddDate(0, 0, -60)
	if !p.IsExpired() {
		t.Fatal("policy should be expired")
	}
}

func TestNewSecretRetentionRegistry_NotNil(t *testing.T) {
	if NewSecretRetentionRegistry() == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestRetentionRegistry_Set_And_Get(t *testing.T) {
	reg := NewSecretRetentionRegistry()
	if err := reg.Set(baseRetention); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get(baseRetention.Mount, baseRetention.Path)
	if !ok {
		t.Fatal("expected policy to be found")
	}
	if got.Duration != baseRetention.Duration {
		t.Errorf("duration mismatch: got %d", got.Duration)
	}
}

func TestRetentionRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretRetentionRegistry()
	_, ok := reg.Get("secret", "missing/path")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestRetentionRegistry_Remove(t *testing.T) {
	reg := NewSecretRetentionRegistry()
	_ = reg.Set(baseRetention)
	reg.Remove(baseRetention.Mount, baseRetention.Path)
	_, ok := reg.Get(baseRetention.Mount, baseRetention.Path)
	if ok {
		t.Fatal("expected policy to be removed")
	}
}

func TestRetentionRegistry_Expired(t *testing.T) {
	reg := NewSecretRetentionRegistry()
	expired := baseRetention
	expired.CreatedAt = time.Now().UTC().AddDate(0, 0, -90)
	_ = reg.Set(expired)
	results := reg.Expired()
	if len(results) != 1 {
		t.Fatalf("expected 1 expired policy, got %d", len(results))
	}
}
