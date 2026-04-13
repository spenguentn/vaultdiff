package vault

import (
	"testing"
)

func TestIsValidVisibilityLevel_Known(t *testing.T) {
	for _, l := range []VisibilityLevel{VisibilityPublic, VisibilityInternal, VisibilityPrivate, VisibilityRestricted} {
		if !IsValidVisibilityLevel(l) {
			t.Errorf("expected %q to be valid", l)
		}
	}
}

func TestIsValidVisibilityLevel_Unknown(t *testing.T) {
	if IsValidVisibilityLevel("ghost") {
		t.Error("expected unknown level to be invalid")
	}
}

func TestSecretVisibility_FullPath(t *testing.T) {
	v := SecretVisibility{Mount: "secret", Path: "app/db", Level: VisibilityPrivate, SetBy: "ops"}
	if got := v.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretVisibility_Validate_Valid(t *testing.T) {
	v := SecretVisibility{Mount: "secret", Path: "app/db", Level: VisibilityInternal, SetBy: "admin"}
	if err := v.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretVisibility_Validate_MissingMount(t *testing.T) {
	v := SecretVisibility{Path: "app/db", Level: VisibilityPublic, SetBy: "admin"}
	if err := v.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretVisibility_Validate_MissingPath(t *testing.T) {
	v := SecretVisibility{Mount: "secret", Level: VisibilityPublic, SetBy: "admin"}
	if err := v.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretVisibility_Validate_InvalidLevel(t *testing.T) {
	v := SecretVisibility{Mount: "secret", Path: "app/db", Level: "ghost", SetBy: "admin"}
	if err := v.Validate(); err == nil {
		t.Error("expected error for invalid level")
	}
}

func TestSecretVisibility_Validate_MissingSetBy(t *testing.T) {
	v := SecretVisibility{Mount: "secret", Path: "app/db", Level: VisibilityRestricted}
	if err := v.Validate(); err == nil {
		t.Error("expected error for missing set_by")
	}
}

func TestNewSecretVisibilityRegistry_NotNil(t *testing.T) {
	if NewSecretVisibilityRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestVisibilityRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretVisibilityRegistry()
	v := SecretVisibility{Mount: "secret", Path: "app/key", Level: VisibilityPrivate, SetBy: "alice"}
	if err := r.Set(v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/key")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.Level != VisibilityPrivate {
		t.Errorf("unexpected level: %s", got.Level)
	}
}

func TestVisibilityRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretVisibilityRegistry()
	_, ok := r.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestVisibilityRegistry_Remove(t *testing.T) {
	r := NewSecretVisibilityRegistry()
	v := SecretVisibility{Mount: "secret", Path: "app/key", Level: VisibilityInternal, SetBy: "bob"}
	_ = r.Set(v)
	r.Remove("secret", "app/key")
	_, ok := r.Get("secret", "app/key")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestVisibilityRegistry_All(t *testing.T) {
	r := NewSecretVisibilityRegistry()
	_ = r.Set(SecretVisibility{Mount: "m", Path: "p1", Level: VisibilityPublic, SetBy: "u"})
	_ = r.Set(SecretVisibility{Mount: "m", Path: "p2", Level: VisibilityPrivate, SetBy: "u"})
	if len(r.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.All()))
	}
}
