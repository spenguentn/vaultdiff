package vault

import (
	"testing"
	"time"
)

func TestIsValidReadinessStatus_Known(t *testing.T) {
	for _, s := range []ReadinessStatus{
		ReadinessReady, ReadinessNotReady, ReadinessPending, ReadinessUnknown,
	} {
		if !IsValidReadinessStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidReadinessStatus_Unknown(t *testing.T) {
	if IsValidReadinessStatus(ReadinessStatus("bogus")) {
		t.Error("expected 'bogus' to be invalid")
	}
}

func TestSecretReadiness_FullPath(t *testing.T) {
	r := SecretReadiness{Mount: "secret", Path: "app/db"}
	if got := r.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected FullPath: %s", got)
	}
}

func TestSecretReadiness_IsReady_True(t *testing.T) {
	r := SecretReadiness{Status: ReadinessReady}
	if !r.IsReady() {
		t.Error("expected IsReady to be true")
	}
}

func TestSecretReadiness_IsReady_False(t *testing.T) {
	r := SecretReadiness{Status: ReadinessNotReady}
	if r.IsReady() {
		t.Error("expected IsReady to be false")
	}
}

func TestSecretReadiness_Validate_Valid(t *testing.T) {
	r := SecretReadiness{
		Mount:     "secret",
		Path:      "app/key",
		Status:    ReadinessReady,
		CheckedBy: "ci-bot",
		CheckedAt: time.Now(),
	}
	if err := r.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretReadiness_Validate_MissingMount(t *testing.T) {
	r := SecretReadiness{Path: "app/key", Status: ReadinessReady, CheckedBy: "ci-bot"}
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretReadiness_Validate_MissingPath(t *testing.T) {
	r := SecretReadiness{Mount: "secret", Status: ReadinessReady, CheckedBy: "ci-bot"}
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretReadiness_Validate_MissingCheckedBy(t *testing.T) {
	r := SecretReadiness{Mount: "secret", Path: "app/key", Status: ReadinessReady}
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing checked_by")
	}
}

func TestSecretReadiness_Validate_InvalidStatus(t *testing.T) {
	r := SecretReadiness{
		Mount:     "secret",
		Path:      "app/key",
		Status:    ReadinessStatus("invalid"),
		CheckedBy: "ci-bot",
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for invalid status")
	}
}
