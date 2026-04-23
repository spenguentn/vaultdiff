package vault

import (
	"testing"
	"time"
)

func TestIsValidTenureStatus_Known(t *testing.T) {
	for _, s := range []TenureStatus{TenureNew, TenureActive, TenureMature, TenureVeteran} {
		if !IsValidTenureStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidTenureStatus_Unknown(t *testing.T) {
	if IsValidTenureStatus("ancient") {
		t.Error("expected 'ancient' to be invalid")
	}
}

func TestSecretTenure_FullPath(t *testing.T) {
	ten := SecretTenure{Mount: "secret", Path: "app/db"}
	if got := ten.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestComputeTenure_New(t *testing.T) {
	created := time.Now().Add(-10 * 24 * time.Hour)
	ten := ComputeTenure("secret", "app/key", created)
	if ten.Status != TenureNew {
		t.Errorf("expected TenureNew, got %s", ten.Status)
	}
}

func TestComputeTenure_Veteran(t *testing.T) {
	created := time.Now().Add(-400 * 24 * time.Hour)
	ten := ComputeTenure("secret", "app/key", created)
	if ten.Status != TenureVeteran {
		t.Errorf("expected TenureVeteran, got %s", ten.Status)
	}
}

func TestSecretTenure_Validate_Valid(t *testing.T) {
	ten := SecretTenure{
		Mount:     "secret",
		Path:      "app/key",
		CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
		Status:    TenureActive,
	}
	if err := ten.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretTenure_Validate_MissingMount(t *testing.T) {
	ten := SecretTenure{Path: "app/key", CreatedAt: time.Now(), Status: TenureNew}
	if err := ten.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretTenure_Validate_ZeroCreatedAt(t *testing.T) {
	ten := SecretTenure{Mount: "secret", Path: "app/key", Status: TenureNew}
	if err := ten.Validate(); err == nil {
		t.Error("expected error for zero created_at")
	}
}

func TestSecretTenure_AgeDays(t *testing.T) {
	ten := SecretTenure{CreatedAt: time.Now().Add(-5 * 24 * time.Hour)}
	if ten.AgeDays() < 4 || ten.AgeDays() > 6 {
		t.Errorf("unexpected age days: %d", ten.AgeDays())
	}
}
