package vault

import (
	"testing"
	"time"
)

func baseAnnotation() SecretAnnotation {
	return SecretAnnotation{
		Mount:     "secret",
		Path:      "app/config",
		Key:       "owner",
		Value:     "platform-team",
		CreatedBy: "alice",
	}
}

func TestSecretAnnotation_FullPath(t *testing.T) {
	a := baseAnnotation()
	if got := a.FullPath(); got != "secret/app/config" {
		t.Errorf("expected secret/app/config, got %s", got)
	}
}

func TestSecretAnnotation_Validate_Valid(t *testing.T) {
	if err := baseAnnotation().Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretAnnotation_Validate_MissingMount(t *testing.T) {
	a := baseAnnotation()
	a.Mount = ""
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretAnnotation_Validate_MissingPath(t *testing.T) {
	a := baseAnnotation()
	a.Path = ""
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretAnnotation_Validate_MissingKey(t *testing.T) {
	a := baseAnnotation()
	a.Key = ""
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestSecretAnnotation_Validate_MissingCreatedBy(t *testing.T) {
	a := baseAnnotation()
	a.CreatedBy = ""
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing created_by")
	}
}

func TestNewSecretAnnotationRegistry_NotNil(t *testing.T) {
	if NewSecretAnnotationRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestAnnotationRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretAnnotationRegistry()
	a := baseAnnotation()
	if err := r.Set(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get(a.Mount, a.Path, a.Key)
	if !ok {
		t.Fatal("expected annotation to be found")
	}
	if got.Value != "platform-team" {
		t.Errorf("expected platform-team, got %s", got.Value)
	}
}

func TestAnnotationRegistry_Set_SetsCreatedAt(t *testing.T) {
	r := NewSecretAnnotationRegistry()
	a := baseAnnotation()
	_ = r.Set(a)
	got, _ := r.Get(a.Mount, a.Path, a.Key)
	if got.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be stamped")
	}
}

func TestAnnotationRegistry_Set_PreservesCreatedAt(t *testing.T) {
	r := NewSecretAnnotationRegistry()
	a := baseAnnotation()
	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a.CreatedAt = fixed
	_ = r.Set(a)
	got, _ := r.Get(a.Mount, a.Path, a.Key)
	if !got.CreatedAt.Equal(fixed) {
		t.Errorf("expected %v, got %v", fixed, got.CreatedAt)
	}
}

func TestAnnotationRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretAnnotationRegistry()
	a := baseAnnotation()
	a.Mount = ""
	if err := r.Set(a); err == nil {
		t.Error("expected validation error")
	}
}

func TestAnnotationRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretAnnotationRegistry()
	_, ok := r.Get("secret", "missing", "key")
	if ok {
		t.Error("expected not found")
	}
}

func TestAnnotationRegistry_Remove(t *testing.T) {
	r := NewSecretAnnotationRegistry()
	a := baseAnnotation()
	_ = r.Set(a)
	r.Remove(a.Mount, a.Path, a.Key)
	_, ok := r.Get(a.Mount, a.Path, a.Key)
	if ok {
		t.Error("expected annotation to be removed")
	}
}

func TestAnnotationRegistry_ListForPath(t *testing.T) {
	r := NewSecretAnnotationRegistry()
	a1 := baseAnnotation()
	a2 := baseAnnotation()
	a2.Key = "team"
	a2.Value = "infra"
	_ = r.Set(a1)
	_ = r.Set(a2)
	list := r.ListForPath("secret", "app/config")
	if len(list) != 2 {
		t.Errorf("expected 2 annotations, got %d", len(list))
	}
}
