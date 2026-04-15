package vault

import (
	"testing"
	"time"
)

var baseRisk = SecretRisk{
	Mount:      "secret",
	Path:       "infra/db",
	Level:      RiskLevelHigh,
	Reason:     "exposed in logs",
	AssessedBy: "alice",
	AssessedAt: time.Now(),
}

func TestIsValidRiskLevel_Known(t *testing.T) {
	for _, lvl := range []RiskLevel{RiskLevelLow, RiskLevelMedium, RiskLevelHigh, RiskLevelCritical} {
		if !IsValidRiskLevel(lvl) {
			t.Errorf("expected %q to be valid", lvl)
		}
	}
}

func TestIsValidRiskLevel_Unknown(t *testing.T) {
	if IsValidRiskLevel("extreme") {
		t.Error("expected 'extreme' to be invalid")
	}
}

func TestSecretRisk_FullPath(t *testing.T) {
	r := baseRisk
	got := r.FullPath()
	want := "secret/infra/db"
	if got != want {
		t.Errorf("FullPath() = %q, want %q", got, want)
	}
}

func TestSecretRisk_Validate_Valid(t *testing.T) {
	r := baseRisk
	if err := r.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretRisk_Validate_MissingMount(t *testing.T) {
	r := baseRisk
	r.Mount = ""
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretRisk_Validate_MissingPath(t *testing.T) {
	r := baseRisk
	r.Path = ""
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretRisk_Validate_InvalidLevel(t *testing.T) {
	r := baseRisk
	r.Level = "unknown"
	if err := r.Validate(); err == nil {
		t.Error("expected error for invalid level")
	}
}

func TestSecretRisk_Validate_MissingAssessedBy(t *testing.T) {
	r := baseRisk
	r.AssessedBy = ""
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing assessed_by")
	}
}
