package vault

import (
	"testing"
)

func baseLabel() SecretLabel {
	return SecretLabel{
		Mount:     "secret",
		Path:      "app/config",
		Key:       "team",
		Value:     "platform",
		CreatedBy: "alice",
	}
}

func TestNewSecretLabelRegistry_NotNil(t *testing.T) {
	r := NewSecretLabelRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestLabelRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretLabelRegistry()
	l := baseLabel()
	if err := r.Set(l); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := r.Get(l.Mount, l.Path, l.Key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Value != "platform" {
		t.Errorf("expected platform, got %s", got.Value)
	}
}

func TestLabelRegistry_Set_SetsCreatedAt(t *testing.T) {
	r := NewSecretLabelRegistry()
	l := baseLabel()
	_ = r.Set(l)
	got, _ := r.Get(l.Mount, l.Path, l.Key)
	if got.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestLabelRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretLabelRegistry()
	l := SecretLabel{Path: "app/config", Key: "team", CreatedBy: "alice"}
	if err := r.Set(l); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestLabelRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretLabelRegistry()
	_, err := r.Get("secret", "app/config", "missing")
	if err == nil {
		t.Fatal("expected not-found error")
	}
}

func TestLabelRegistry_List(t *testing.T) {
	r := NewSecretLabelRegistry()
	l1 := baseLabel()
	l2 := baseLabel()
	l2.Key = "env"
	l2.Value = "prod"
	_ = r.Set(l1)
	_ = r.Set(l2)
	results := r.List("secret", "app/config")
	if len(results) != 2 {
		t.Errorf("expected 2 labels, got %d", len(results))
	}
}

func TestLabelRegistry_Remove(t *testing.T) {
	r := NewSecretLabelRegistry()
	l := baseLabel()
	_ = r.Set(l)
	if err := r.Remove(l.Mount, l.Path, l.Key); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := r.Get(l.Mount, l.Path, l.Key); err == nil {
		t.Fatal("expected not-found after remove")
	}
}

func TestLabelRegistry_Remove_NotFound(t *testing.T) {
	r := NewSecretLabelRegistry()
	if err := r.Remove("secret", "app/config", "ghost"); err == nil {
		t.Fatal("expected error for missing label")
	}
}
