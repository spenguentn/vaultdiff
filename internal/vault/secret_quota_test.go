package vault

import (
	"testing"
	"time"
)

var baseQuota = SecretQuota{
	Mount:      "secret",
	Prefix:     "app/config",
	Scope:      QuotaScopeMount,
	MaxReads:   10,
	WindowSize: time.Minute,
}

func TestSecretQuota_FullPath(t *testing.T) {
	q := baseQuota
	want := "secret/app/config"
	if got := q.FullPath(); got != want {
		t.Errorf("FullPath() = %q, want %q", got, want)
	}
}

func TestSecretQuota_Validate_Valid(t *testing.T) {
	if err := baseQuota.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretQuota_Validate_MissingMount(t *testing.T) {
	q := baseQuota
	q.Mount = ""
	if err := q.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretQuota_Validate_ZeroMaxReads(t *testing.T) {
	q := baseQuota
	q.MaxReads = 0
	if err := q.Validate(); err == nil {
		t.Error("expected error for zero max_reads")
	}
}

func TestSecretQuota_Validate_ZeroWindow(t *testing.T) {
	q := baseQuota
	q.WindowSize = 0
	if err := q.Validate(); err == nil {
		t.Error("expected error for zero window_size")
	}
}

func TestSecretQuota_Validate_MissingScope(t *testing.T) {
	q := baseQuota
	q.Scope = ""
	if err := q.Validate(); err == nil {
		t.Error("expected error for missing scope")
	}
}

func TestNewSecretQuotaRegistry_NotNil(t *testing.T) {
	if r := NewSecretQuotaRegistry(); r == nil {
		t.Error("expected non-nil registry")
	}
}

func TestQuotaRegistry_Register_And_Get(t *testing.T) {
	r := NewSecretQuotaRegistry()
	if err := r.Register(baseQuota); err != nil {
		t.Fatalf("Register() error: %v", err)
	}
	got, ok := r.Get(baseQuota.Mount, baseQuota.Prefix)
	if !ok {
		t.Fatal("Get() returned not found")
	}
	if got.MaxReads != baseQuota.MaxReads {
		t.Errorf("MaxReads = %d, want %d", got.MaxReads, baseQuota.MaxReads)
	}
}

func TestQuotaRegistry_Register_Invalid(t *testing.T) {
	r := NewSecretQuotaRegistry()
	q := baseQuota
	q.Mount = ""
	if err := r.Register(q); err == nil {
		t.Error("expected error for invalid quota")
	}
}

func TestQuotaRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretQuotaRegistry()
	_, ok := r.Get("missing", "path")
	if ok {
		t.Error("expected not found")
	}
}

func TestQuotaRegistry_Remove(t *testing.T) {
	r := NewSecretQuotaRegistry()
	_ = r.Register(baseQuota)
	r.Remove(baseQuota.Mount, baseQuota.Prefix)
	_, ok := r.Get(baseQuota.Mount, baseQuota.Prefix)
	if ok {
		t.Error("expected quota to be removed")
	}
}

func TestQuotaRegistry_Allow_WithinLimit(t *testing.T) {
	r := NewSecretQuotaRegistry()
	q := baseQuota
	q.MaxReads = 3
	q.WindowSize = time.Minute
	_ = r.Register(q)
	for i := 0; i < 3; i++ {
		if !r.Allow(q.Mount, q.Prefix) {
			t.Errorf("Allow() returned false on attempt %d", i+1)
		}
	}
}

func TestQuotaRegistry_Allow_ExceedsLimit(t *testing.T) {
	r := NewSecretQuotaRegistry()
	q := baseQuota
	q.MaxReads = 2
	q.WindowSize = time.Minute
	_ = r.Register(q)
	r.Allow(q.Mount, q.Prefix)
	r.Allow(q.Mount, q.Prefix)
	if r.Allow(q.Mount, q.Prefix) {
		t.Error("Allow() should return false when quota exceeded")
	}
}

func TestQuotaRegistry_Allow_NoQuota_Permits(t *testing.T) {
	r := NewSecretQuotaRegistry()
	if !r.Allow("any", "path") {
		t.Error("Allow() should return true when no quota registered")
	}
}
