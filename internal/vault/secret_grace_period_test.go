package vault

import (
	"testing"
	"time"
)

func baseGracePeriod() SecretGracePeriod {
	return SecretGracePeriod{
		Mount:     "secret",
		Path:      "app/db-password",
		Status:    GracePeriodActive,
		StartsAt:  time.Now().Add(-time.Hour),
		Duration:  48 * time.Hour,
		GrantedBy: "ops-team",
		Reason:    "scheduled maintenance window",
	}
}

func TestIsValidGracePeriodStatus_Known(t *testing.T) {
	for _, s := range []GracePeriodStatus{GracePeriodActive, GracePeriodExpired, GracePeriodPending} {
		if !IsValidGracePeriodStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidGracePeriodStatus_Unknown(t *testing.T) {
	if IsValidGracePeriodStatus("unknown") {
		t.Error("expected unknown status to be invalid")
	}
}

func TestSecretGracePeriod_FullPath(t *testing.T) {
	g := baseGracePeriod()
	if got := g.FullPath(); got != "secret/app/db-password" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretGracePeriod_IsExpired_False(t *testing.T) {
	g := baseGracePeriod()
	if g.IsExpired() {
		t.Error("expected grace period to not be expired")
	}
}

func TestSecretGracePeriod_IsExpired_True(t *testing.T) {
	g := baseGracePeriod()
	g.StartsAt = time.Now().Add(-72 * time.Hour)
	g.Duration = time.Hour
	if !g.IsExpired() {
		t.Error("expected grace period to be expired")
	}
}

func TestSecretGracePeriod_Validate_Valid(t *testing.T) {
	if err := baseGracePeriod().Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestSecretGracePeriod_Validate_MissingMount(t *testing.T) {
	g := baseGracePeriod()
	g.Mount = ""
	if err := g.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretGracePeriod_Validate_MissingGrantedBy(t *testing.T) {
	g := baseGracePeriod()
	g.GrantedBy = ""
	if err := g.Validate(); err == nil {
		t.Error("expected error for missing granted_by")
	}
}

func TestSecretGracePeriod_Validate_ZeroDuration(t *testing.T) {
	g := baseGracePeriod()
	g.Duration = 0
	if err := g.Validate(); err == nil {
		t.Error("expected error for zero duration")
	}
}

func TestSecretGracePeriod_Validate_InvalidStatus(t *testing.T) {
	g := baseGracePeriod()
	g.Status = "invalid"
	if err := g.Validate(); err == nil {
		t.Error("expected error for invalid status")
	}
}
