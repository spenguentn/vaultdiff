package vault

import (
	"testing"
	"time"
)

var baseBlame = BlameEntry{
	Mount:     "secret",
	Path:      "app/config",
	Version:   2,
	ChangedBy: "alice",
	ChangedAt: time.Now(),
	Operation: "write",
}

func TestBlameEntry_FullPath(t *testing.T) {
	if got := baseBlame.FullPath(); got != "secret/app/config" {
		t.Fatalf("expected secret/app/config, got %s", got)
	}
}

func TestBlameEntry_Validate_Valid(t *testing.T) {
	if err := baseBlame.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBlameEntry_Validate_MissingMount(t *testing.T) {
	e := baseBlame
	e.Mount = ""
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestBlameEntry_Validate_MissingChangedBy(t *testing.T) {
	e := baseBlame
	e.ChangedBy = ""
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for missing changed_by")
	}
}

func TestBlameEntry_Validate_ZeroVersion(t *testing.T) {
	e := baseBlame
	e.Version = 0
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for zero version")
	}
}

func TestIsValidOperation_Known(t *testing.T) {
	for _, op := range []string{"write", "delete", "restore"} {
		if !IsValidOperation(op) {
			t.Fatalf("expected %s to be valid", op)
		}
	}
}

func TestIsValidOperation_Unknown(t *testing.T) {
	if IsValidOperation("rotate") {
		t.Fatal("expected rotate to be invalid")
	}
}

func TestNewSecretBlameRegistry_NotNil(t *testing.T) {
	if NewSecretBlameRegistry() == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestBlameRegistry_Record_And_Get(t *testing.T) {
	r := NewSecretBlameRegistry()
	if err := r.Record(baseBlame); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get(baseBlame.Mount, baseBlame.Path, baseBlame.Version)
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.ChangedBy != baseBlame.ChangedBy {
		t.Fatalf("expected %s, got %s", baseBlame.ChangedBy, got.ChangedBy)
	}
}

func TestBlameRegistry_Record_Invalid(t *testing.T) {
	r := NewSecretBlameRegistry()
	e := baseBlame
	e.Path = ""
	if err := r.Record(e); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestBlameRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretBlameRegistry()
	_, ok := r.Get("secret", "missing", 1)
	if ok {
		t.Fatal("expected not found")
	}
}

func TestBlameRegistry_Latest(t *testing.T) {
	r := NewSecretBlameRegistry()
	v1 := baseBlame
	v1.Version = 1
	v2 := baseBlame
	v2.Version = 2
	_ = r.Record(v1)
	_ = r.Record(v2)
	got, ok := r.Latest(baseBlame.Mount, baseBlame.Path)
	if !ok || got.Version != 2 {
		t.Fatalf("expected version 2, got %d", got.Version)
	}
}

func TestBlameRegistry_Remove(t *testing.T) {
	r := NewSecretBlameRegistry()
	_ = r.Record(baseBlame)
	r.Remove(baseBlame.Mount, baseBlame.Path, baseBlame.Version)
	_, ok := r.Get(baseBlame.Mount, baseBlame.Path, baseBlame.Version)
	if ok {
		t.Fatal("expected entry to be removed")
	}
}
