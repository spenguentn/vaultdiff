package vault

import (
	"testing"
	"time"
)

func baseDeprecation() SecretDeprecation {
	return SecretDeprecation{
		Mount:        "secret",
		Path:         "myapp/db-password",
		Status:       DeprecationStatusDeprecated,
		Reason:       "rotated to new path",
		DeprecatedBy: "ops-team",
		DeprecatedAt: time.Now(),
	}
}

func TestIsValidDeprecationStatus_Known(t *testing.T) {
	for _, s := range []DeprecationStatus{
		DeprecationStatusActive,
		DeprecationStatusDeprecated,
		DeprecationStatusSunset,
	} {
		if !IsValidDeprecationStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidDeprecationStatus_Unknown(t *testing.T) {
	if IsValidDeprecationStatus("unknown") {
		t.Error("expected unknown status to be invalid")
	}
}

func TestSecretDeprecation_FullPath(t *testing.T) {
	d := baseDeprecation()
	if got := d.FullPath(); got != "secret/myapp/db-password" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretDeprecation_IsSunset_False(t *testing.T) {
	d := baseDeprecation()
	future := time.Now().Add(24 * time.Hour)
	d.SunsetAt = &future
	if d.IsSunset() {
		t.Error("expected IsSunset to be false for future sunset")
	}
}

func TestSecretDeprecation_IsSunset_True(t *testing.T) {
	d := baseDeprecation()
	past := time.Now().Add(-time.Hour)
	d.SunsetAt = &past
	if !d.IsSunset() {
		t.Error("expected IsSunset to be true for past sunset")
	}
}

func TestSecretDeprecation_IsSunset_NoExpiry(t *testing.T) {
	d := baseDeprecation()
	if d.IsSunset() {
		t.Error("expected IsSunset to be false when SunsetAt is nil")
	}
}

func TestSecretDeprecation_Validate_Valid(t *testing.T) {
	if err := baseDeprecation().Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretDeprecation_Validate_MissingMount(t *testing.T) {
	d := baseDeprecation()
	d.Mount = ""
	if err := d.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretDeprecation_Validate_MissingDeprecatedBy(t *testing.T) {
	d := baseDeprecation()
	d.DeprecatedBy = ""
	if err := d.Validate(); err == nil {
		t.Error("expected error for missing deprecated_by")
	}
}

func TestSecretDeprecation_Validate_InvalidStatus(t *testing.T) {
	d := baseDeprecation()
	d.Status = "bogus"
	if err := d.Validate(); err == nil {
		t.Error("expected error for invalid status")
	}
}
