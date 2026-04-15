package vault

import (
	"testing"
	"time"
)

func TestIsValidMaturityLevel_Known(t *testing.T) {
	for _, l := range []MaturityLevel{MaturityDraft, MaturityCandidate, MaturityStable, MaturityDeprecated} {
		if !IsValidMaturityLevel(l) {
			t.Errorf("expected %q to be valid", l)
		}
	}
}

func TestIsValidMaturityLevel_Unknown(t *testing.T) {
	if IsValidMaturityLevel(MaturityLevel("unknown")) {
		t.Error("expected unknown level to be invalid")
	}
}

func TestSecretMaturity_FullPath(t *testing.T) {
	m := SecretMaturity{Mount: "secret", Path: "app/db"}
	if got := m.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretMaturity_Validate_Valid(t *testing.T) {
	m := SecretMaturity{Mount: "secret", Path: "app/db", Level: MaturityStable, AssignedBy: "alice"}
	if err := m.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretMaturity_Validate_MissingMount(t *testing.T) {
	m := SecretMaturity{Path: "app/db", Level: MaturityStable, AssignedBy: "alice"}
	if err := m.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretMaturity_Validate_MissingPath(t *testing.T) {
	m := SecretMaturity{Mount: "secret", Level: MaturityStable, AssignedBy: "alice"}
	if err := m.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretMaturity_Validate_UnknownLevel(t *testing.T) {
	m := SecretMaturity{Mount: "secret", Path: "app/db", Level: "experimental", AssignedBy: "alice"}
	if err := m.Validate(); err == nil {
		t.Error("expected error for unknown level")
	}
}

func TestMaturityRegistry_Set_And_Get(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	m := SecretMaturity{Mount: "secret", Path: "app/db", Level: MaturityStable, AssignedBy: "alice"}
	if err := reg.Set(m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := reg.Get("secret", "app/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Level != MaturityStable {
		t.Errorf("expected stable, got %s", got.Level)
	}
}

func TestMaturityRegistry_Set_SetsAssignedAt(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	m := SecretMaturity{Mount: "secret", Path: "key", Level: MaturityDraft, AssignedBy: "bob"}
	_ = reg.Set(m)
	got, _ := reg.Get("secret", "key")
	if got.AssignedAt.IsZero() {
		t.Error("expected AssignedAt to be set automatically")
	}
}

func TestMaturityRegistry_Set_Invalid(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	if err := reg.Set(SecretMaturity{}); err == nil {
		t.Error("expected validation error")
	}
}

func TestMaturityRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	if _, err := reg.Get("secret", "missing"); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestMaturityRegistry_Remove(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	m := SecretMaturity{Mount: "secret", Path: "key", Level: MaturityCandidate, AssignedBy: "carol", AssignedAt: time.Now()}
	_ = reg.Set(m)
	reg.Remove("secret", "key")
	if _, err := reg.Get("secret", "key"); err == nil {
		t.Error("expected error after removal")
	}
}

func TestMaturityRegistry_All(t *testing.T) {
	reg := NewSecretMaturityRegistry()
	for _, path := range []string{"a", "b", "c"} {
		_ = reg.Set(SecretMaturity{Mount: "m", Path: path, Level: MaturityStable, AssignedBy: "x", AssignedAt: time.Now()})
	}
	if len(reg.All()) != 3 {
		t.Errorf("expected 3 entries, got %d", len(reg.All()))
	}
}
