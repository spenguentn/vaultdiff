package vault

import "testing"

func baseCategory() SecretCategory {
	return SecretCategory{
		Mount:    "secret",
		Path:     "myapp/db",
		Category: "credentials",
		SetBy:    "admin",
	}
}

func TestIsValidCategory_Known(t *testing.T) {
	for _, c := range ValidCategories {
		if !IsValidCategory(c) {
			t.Errorf("expected %q to be valid", c)
		}
	}
}

func TestIsValidCategory_Unknown(t *testing.T) {
	if IsValidCategory("bogus") {
		t.Error("expected 'bogus' to be invalid")
	}
}

func TestSecretCategory_FullPath(t *testing.T) {
	sc := baseCategory()
	if got := sc.FullPath(); got != "secret/myapp/db" {
		t.Errorf("unexpected FullPath: %s", got)
	}
}

func TestSecretCategory_Validate_Valid(t *testing.T) {
	if err := baseCategory().Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretCategory_Validate_MissingMount(t *testing.T) {
	sc := baseCategory()
	sc.Mount = ""
	if err := sc.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretCategory_Validate_MissingPath(t *testing.T) {
	sc := baseCategory()
	sc.Path = ""
	if err := sc.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretCategory_Validate_UnknownCategory(t *testing.T) {
	sc := baseCategory()
	sc.Category = "unknown-cat"
	if err := sc.Validate(); err == nil {
		t.Error("expected error for unknown category")
	}
}

func TestSecretCategory_Validate_MissingSetBy(t *testing.T) {
	sc := baseCategory()
	sc.SetBy = ""
	if err := sc.Validate(); err == nil {
		t.Error("expected error for missing set_by")
	}
}

func TestNewSecretCategoryRegistry_NotNil(t *testing.T) {
	if NewSecretCategoryRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestCategoryRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretCategoryRegistry()
	sc := baseCategory()
	if err := r.Set(sc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get(sc.Mount, sc.Path)
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.Category != sc.Category {
		t.Errorf("got category %q, want %q", got.Category, sc.Category)
	}
}

func TestCategoryRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretCategoryRegistry()
	sc := baseCategory()
	sc.Mount = ""
	if err := r.Set(sc); err == nil {
		t.Error("expected error for invalid entry")
	}
}

func TestCategoryRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretCategoryRegistry()
	_, ok := r.Get("secret", "missing/path")
	if ok {
		t.Error("expected not found")
	}
}

func TestCategoryRegistry_Remove(t *testing.T) {
	r := NewSecretCategoryRegistry()
	sc := baseCategory()
	_ = r.Set(sc)
	r.Remove(sc.Mount, sc.Path)
	if r.Len() != 0 {
		t.Errorf("expected 0 entries after remove, got %d", r.Len())
	}
}
