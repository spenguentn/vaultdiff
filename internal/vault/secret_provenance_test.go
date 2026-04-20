package vault

import (
	"testing"
	"time"
)

func TestIsValidProvenanceSource_Known(t *testing.T) {
	known := []string{"manual", "automated", "imported", "generated", "inherited"}
	for _, s := range known {
		if !IsValidProvenanceSource(s) {
			t.Errorf("expected %q to be a valid provenance source", s)
		}
	}
}

func TestIsValidProvenanceSource_Unknown(t *testing.T) {
	if IsValidProvenanceSource("unknown-source") {
		t.Error("expected 'unknown-source' to be invalid")
	}
}

func TestSecretProvenance_FullPath(t *testing.T) {
	p := SecretProvenance{
		Mount: "secret",
		Path:  "db/password",
	}
	want := "secret/db/password"
	if got := p.FullPath(); got != want {
		t.Errorf("FullPath() = %q; want %q", got, want)
	}
}

func TestSecretProvenance_Validate_Valid(t *testing.T) {
	p := SecretProvenance{
		Mount:     "secret",
		Path:      "db/password",
		Source:    "manual",
		CreatedBy: "alice",
		CreatedAt: time.Now(),
	}
	if err := p.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretProvenance_Validate_MissingMount(t *testing.T) {
	p := SecretProvenance{
		Path:      "db/password",
		Source:    "manual",
		CreatedBy: "alice",
		CreatedAt: time.Now(),
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretProvenance_Validate_MissingPath(t *testing.T) {
	p := SecretProvenance{
		Mount:     "secret",
		Source:    "manual",
		CreatedBy: "alice",
		CreatedAt: time.Now(),
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretProvenance_Validate_InvalidSource(t *testing.T) {
	p := SecretProvenance{
		Mount:     "secret",
		Path:      "db/password",
		Source:    "unknown",
		CreatedBy: "alice",
		CreatedAt: time.Now(),
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for invalid source")
	}
}

func TestSecretProvenance_Validate_MissingCreatedBy(t *testing.T) {
	p := SecretProvenance{
		Mount:     "secret",
		Path:      "db/password",
		Source:    "automated",
		CreatedAt: time.Now(),
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing created_by")
	}
}

func TestProvenanceRegistry_Set_And_Get(t *testing.T) {
	reg := NewSecretProvenanceRegistry()
	p := SecretProvenance{
		Mount:     "secret",
		Path:      "db/password",
		Source:    "manual",
		CreatedBy: "alice",
		CreatedAt: time.Now(),
	}
	if err := reg.Set(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("secret", "db/password")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.CreatedBy != "alice" {
		t.Errorf("CreatedBy = %q; want %q", got.CreatedBy, "alice")
	}
}

func TestProvenanceRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretProvenanceRegistry()
	_, ok := reg.Get("secret", "missing/path")
	if ok {
		t.Error("expected not found")
	}
}

func TestProvenanceRegistry_Remove(t *testing.T) {
	reg := NewSecretProvenanceRegistry()
	p := SecretProvenance{
		Mount:     "secret",
		Path:      "db/password",
		Source:    "imported",
		CreatedBy: "bob",
		CreatedAt: time.Now(),
	}
	_ = reg.Set(p)
	reg.Remove("secret", "db/password")
	_, ok := reg.Get("secret", "db/password")
	if ok {
		t.Error("expected entry to be removed")
	}
}
