package vault

import (
	"testing"
	"time"
)

func baseStatus() SecretStatus {
	return SecretStatus{
		Mount:     "secret",
		Path:      "app/db-password",
		Stage:     StatusActive,
		Reason:    "initial creation",
		UpdatedBy: "admin",
		UpdatedAt: time.Now(),
	}
}

func TestIsValidStatus_Known(t *testing.T) {
	for _, s := range []SecretStatusStage{StatusActive, StatusDeprecated, StatusPendingRotation, StatusExpired, StatusRevoked} {
		if !IsValidStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidStatus_Unknown(t *testing.T) {
	if IsValidStatus("unknown-stage") {
		t.Error("expected unknown stage to be invalid")
	}
}

func TestSecretStatus_FullPath(t *testing.T) {
	s := baseStatus()
	if got := s.FullPath(); got != "secret/app/db-password" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretStatus_Validate_Valid(t *testing.T) {
	if err := baseStatus().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretStatus_Validate_MissingMount(t *testing.T) {
	s := baseStatus()
	s.Mount = ""
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretStatus_Validate_MissingPath(t *testing.T) {
	s := baseStatus()
	s.Path = ""
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretStatus_Validate_MissingUpdatedBy(t *testing.T) {
	s := baseStatus()
	s.UpdatedBy = ""
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing updated_by")
	}
}

func TestSecretStatus_Validate_InvalidStage(t *testing.T) {
	s := baseStatus()
	s.Stage = "bad-stage"
	if err := s.Validate(); err == nil {
		t.Error("expected error for invalid stage")
	}
}
