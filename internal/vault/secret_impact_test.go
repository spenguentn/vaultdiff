package vault

import (
	"testing"
)

var baseImpact = SecretImpact{
	Mount:       "secret",
	Path:        "app/db",
	Level:       ImpactHigh,
	AssessedBy:  "alice",
	Justification: "production database credentials",
}

func TestIsValidImpactLevel_Known(t *testing.T) {
	for _, l := range []ImpactLevel{ImpactLow, ImpactMedium, ImpactHigh, ImpactCritical} {
		if !IsValidImpactLevel(l) {
			t.Errorf("expected %q to be valid", l)
		}
	}
}

func TestIsValidImpactLevel_Unknown(t *testing.T) {
	if IsValidImpactLevel(ImpactLevel("extreme")) {
		t.Error("expected unknown level to be invalid")
	}
}

func TestSecretImpact_FullPath(t *testing.T) {
	if got := baseImpact.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretImpact_Validate_Valid(t *testing.T) {
	if err := baseImpact.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretImpact_Validate_MissingMount(t *testing.T) {
	e := baseImpact
	e.Mount = ""
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretImpact_Validate_MissingPath(t *testing.T) {
	e := baseImpact
	e.Path = ""
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretImpact_Validate_MissingAssessedBy(t *testing.T) {
	e := baseImpact
	e.AssessedBy = ""
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing assessed_by")
	}
}

func TestSecretImpact_Validate_UnknownLevel(t *testing.T) {
	e := baseImpact
	e.Level = ImpactLevel("none")
	if err := e.Validate(); err == nil {
		t.Error("expected error for unknown level")
	}
}

func TestNewSecretImpactRegistry_NotNil(t *testing.T) {
	if NewSecretImpactRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestImpactRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretImpactRegistry()
	if err := r.Set(baseImpact); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get(baseImpact.Mount, baseImpact.Path)
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.Level != ImpactHigh {
		t.Errorf("unexpected level: %s", got.Level)
	}
}

func TestImpactRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretImpactRegistry()
	bad := baseImpact
	bad.Mount = ""
	if err := r.Set(bad); err == nil {
		t.Error("expected validation error")
	}
}

func TestImpactRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretImpactRegistry()
	_, ok := r.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestImpactRegistry_Remove(t *testing.T) {
	r := NewSecretImpactRegistry()
	_ = r.Set(baseImpact)
	r.Remove(baseImpact.Mount, baseImpact.Path)
	_, ok := r.Get(baseImpact.Mount, baseImpact.Path)
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestImpactRegistry_All(t *testing.T) {
	r := NewSecretImpactRegistry()
	_ = r.Set(baseImpact)
	if len(r.All()) != 1 {
		t.Errorf("expected 1 entry, got %d", len(r.All()))
	}
}
