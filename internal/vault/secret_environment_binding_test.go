package vault

import (
	"testing"
	"time"
)

var baseBinding = EnvironmentBinding{
	Mount:       "secret",
	Path:        "app/db",
	Environment: "production",
	BoundBy:     "alice",
}

func TestEnvironmentBinding_FullPath(t *testing.T) {
	b := baseBinding
	if got := b.FullPath(); got != "secret/app/db" {
		t.Errorf("expected secret/app/db, got %s", got)
	}
}

func TestEnvironmentBinding_Validate_Valid(t *testing.T) {
	if err := baseBinding.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEnvironmentBinding_Validate_MissingMount(t *testing.T) {
	b := baseBinding
	b.Mount = ""
	if err := b.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestEnvironmentBinding_Validate_MissingPath(t *testing.T) {
	b := baseBinding
	b.Path = ""
	if err := b.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestEnvironmentBinding_Validate_MissingEnvironment(t *testing.T) {
	b := baseBinding
	b.Environment = ""
	if err := b.Validate(); err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestEnvironmentBinding_Validate_MissingBoundBy(t *testing.T) {
	b := baseBinding
	b.BoundBy = ""
	if err := b.Validate(); err == nil {
		t.Error("expected error for missing bound_by")
	}
}

func TestNewEnvironmentBindingRegistry_NotNil(t *testing.T) {
	if r := NewEnvironmentBindingRegistry(); r == nil {
		t.Error("expected non-nil registry")
	}
}

func TestBindingRegistry_Bind_And_Get(t *testing.T) {
	r := NewEnvironmentBindingRegistry()
	if err := r.Bind(baseBinding); err != nil {
		t.Fatalf("unexpected bind error: %v", err)
	}
	got, ok := r.Get("secret", "app/db", "production")
	if !ok {
		t.Fatal("expected binding to be found")
	}
	if got.BoundBy != "alice" {
		t.Errorf("expected alice, got %s", got.BoundBy)
	}
}

func TestBindingRegistry_Bind_SetsTimestamp(t *testing.T) {
	r := NewEnvironmentBindingRegistry()
	before := time.Now().UTC()
	_ = r.Bind(baseBinding)
	got, _ := r.Get("secret", "app/db", "production")
	if got.BoundAt.Before(before) {
		t.Error("expected BoundAt to be set to current time")
	}
}

func TestBindingRegistry_Bind_Invalid(t *testing.T) {
	r := NewEnvironmentBindingRegistry()
	b := baseBinding
	b.Mount = ""
	if err := r.Bind(b); err == nil {
		t.Error("expected error for invalid binding")
	}
}

func TestBindingRegistry_Get_NotFound(t *testing.T) {
	r := NewEnvironmentBindingRegistry()
	_, ok := r.Get("secret", "missing", "staging")
	if ok {
		t.Error("expected not found")
	}
}

func TestBindingRegistry_Remove(t *testing.T) {
	r := NewEnvironmentBindingRegistry()
	_ = r.Bind(baseBinding)
	r.Remove("secret", "app/db", "production")
	_, ok := r.Get("secret", "app/db", "production")
	if ok {
		t.Error("expected binding to be removed")
	}
}

func TestBindingRegistry_All(t *testing.T) {
	r := NewEnvironmentBindingRegistry()
	_ = r.Bind(baseBinding)
	b2 := baseBinding
	b2.Environment = "staging"
	_ = r.Bind(b2)
	if len(r.All()) != 2 {
		t.Errorf("expected 2 bindings, got %d", len(r.All()))
	}
}
