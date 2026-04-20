package vault

import (
	"testing"
	"time"
)

func TestIsValidConfidenceLevel_Known(t *testing.T) {
	for _, lvl := range []ConfidenceLevel{ConfidenceHigh, ConfidenceMedium, ConfidenceLow, ConfidenceUnknown} {
		if !IsValidConfidenceLevel(lvl) {
			t.Errorf("expected %q to be valid", lvl)
		}
	}
}

func TestIsValidConfidenceLevel_Unknown(t *testing.T) {
	if IsValidConfidenceLevel(ConfidenceLevel("bogus")) {
		t.Error("expected bogus level to be invalid")
	}
}

func TestSecretConfidence_FullPath(t *testing.T) {
	sc := SecretConfidence{Mount: "/secret/", Path: "/db/password"}
	if got := sc.FullPath(); got != "secret/db/password" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretConfidence_Validate_Valid(t *testing.T) {
	sc := SecretConfidence{
		Mount:      "secret",
		Path:       "db/password",
		Level:      ConfidenceHigh,
		AssignedBy: "alice",
	}
	if err := sc.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretConfidence_Validate_MissingMount(t *testing.T) {
	sc := SecretConfidence{Path: "db/password", Level: ConfidenceHigh, AssignedBy: "alice"}
	if err := sc.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretConfidence_Validate_MissingPath(t *testing.T) {
	sc := SecretConfidence{Mount: "secret", Level: ConfidenceMedium, AssignedBy: "bob"}
	if err := sc.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretConfidence_Validate_InvalidLevel(t *testing.T) {
	sc := SecretConfidence{Mount: "secret", Path: "db/password", Level: "extreme", AssignedBy: "carol"}
	if err := sc.Validate(); err == nil {
		t.Error("expected error for invalid level")
	}
}

func TestSecretConfidence_Validate_MissingAssignedBy(t *testing.T) {
	sc := SecretConfidence{Mount: "secret", Path: "db/password", Level: ConfidenceLow}
	if err := sc.Validate(); err == nil {
		t.Error("expected error for missing assigned_by")
	}
}

func TestNewSecretConfidenceRegistry_NotNil(t *testing.T) {
	if NewSecretConfidenceRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestConfidenceRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretConfidenceRegistry()
	sc := SecretConfidence{Mount: "secret", Path: "app/key", Level: ConfidenceHigh, AssignedBy: "alice"}
	if err := r.Set(sc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/key")
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.Level != ConfidenceHigh {
		t.Errorf("expected high, got %s", got.Level)
	}
}

func TestConfidenceRegistry_Set_SetsAssignedAt(t *testing.T) {
	r := NewSecretConfidenceRegistry()
	before := time.Now().UTC()
	sc := SecretConfidence{Mount: "secret", Path: "app/key", Level: ConfidenceMedium, AssignedBy: "bob"}
	_ = r.Set(sc)
	got, _ := r.Get("secret", "app/key")
	if got.AssignedAt.Before(before) {
		t.Error("expected AssignedAt to be set to now")
	}
}

func TestConfidenceRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretConfidenceRegistry()
	sc := SecretConfidence{Mount: "", Path: "app/key", Level: ConfidenceLow, AssignedBy: "carol"}
	if err := r.Set(sc); err == nil {
		t.Error("expected validation error")
	}
}

func TestConfidenceRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretConfidenceRegistry()
	_, ok := r.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestConfidenceRegistry_Remove(t *testing.T) {
	r := NewSecretConfidenceRegistry()
	sc := SecretConfidence{Mount: "secret", Path: "app/key", Level: ConfidenceHigh, AssignedBy: "dave"}
	_ = r.Set(sc)
	r.Remove("secret", "app/key")
	_, ok := r.Get("secret", "app/key")
	if ok {
		t.Error("expected record to be removed")
	}
}

func TestConfidenceRegistry_All(t *testing.T) {
	r := NewSecretConfidenceRegistry()
	_ = r.Set(SecretConfidence{Mount: "secret", Path: "a", Level: ConfidenceHigh, AssignedBy: "x"})
	_ = r.Set(SecretConfidence{Mount: "secret", Path: "b", Level: ConfidenceLow, AssignedBy: "y"})
	if len(r.All()) != 2 {
		t.Errorf("expected 2 records, got %d", len(r.All()))
	}
}
