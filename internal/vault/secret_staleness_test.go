package vault

import (
	"testing"
	"time"
)

func TestIsValidStalenessLevel_Known(t *testing.T) {
	for _, l := range []StalenessLevel{
		StalenessLevelFresh, StalenessLevelWarning, StalenessLevelStale, StalenessLevelCritical,
	} {
		if !IsValidStalenessLevel(l) {
			t.Errorf("expected %q to be valid", l)
		}
	}
}

func TestIsValidStalenessLevel_Unknown(t *testing.T) {
	if IsValidStalenessLevel(StalenessLevel("ancient")) {
		t.Error("expected unknown level to be invalid")
	}
}

func TestSecretStaleness_FullPath(t *testing.T) {
	s := SecretStaleness{Mount: "secret", Path: "app/db"}
	if got := s.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretStaleness_AgeDays(t *testing.T) {
	now := time.Now().UTC()
	s := SecretStaleness{
		LastUpdated: now.Add(-72 * time.Hour),
		EvaluatedAt: now,
	}
	if s.AgeDays() != 3 {
		t.Errorf("expected age 3, got %d", s.AgeDays())
	}
}

func TestSecretStaleness_Validate_Valid(t *testing.T) {
	s := SecretStaleness{
		Mount:       "secret",
		Path:        "app/key",
		Level:       StalenessLevelFresh,
		LastUpdated: time.Now(),
	}
	if err := s.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretStaleness_Validate_MissingMount(t *testing.T) {
	s := SecretStaleness{Path: "app/key", Level: StalenessLevelFresh, LastUpdated: time.Now()}
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretStaleness_Validate_InvalidLevel(t *testing.T) {
	s := SecretStaleness{
		Mount:       "secret",
		Path:        "app/key",
		Level:       StalenessLevel("unknown"),
		LastUpdated: time.Now(),
	}
	if err := s.Validate(); err == nil {
		t.Error("expected error for invalid level")
	}
}

func TestComputeStaleness_Fresh(t *testing.T) {
	lastUpdated := time.Now().Add(-24 * time.Hour)
	result, err := ComputeStaleness("secret", "app/key", lastUpdated, 7, 30, 90)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Level != StalenessLevelFresh {
		t.Errorf("expected fresh, got %s", result.Level)
	}
}

func TestComputeStaleness_Critical(t *testing.T) {
	lastUpdated := time.Now().Add(-100 * 24 * time.Hour)
	result, err := ComputeStaleness("secret", "app/key", lastUpdated, 7, 30, 90)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Level != StalenessLevelCritical {
		t.Errorf("expected critical, got %s", result.Level)
	}
}

func TestComputeStaleness_MissingMount(t *testing.T) {
	_, err := ComputeStaleness("", "app/key", time.Now(), 7, 30, 90)
	if err == nil {
		t.Error("expected error for missing mount")
	}
}
