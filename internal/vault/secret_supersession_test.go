package vault

import (
	"testing"
	"time"
)

func TestIsValidSupersessionReason_Known(t *testing.T) {
	known := []SupersessionReason{
		SupersessionReasonRotated,
		SupersessionReasonMigrated,
		SupersessionReasonDeprecated,
		SupersessionReasonReplaced,
	}
	for _, r := range known {
		if !IsValidSupersessionReason(r) {
			t.Errorf("expected %q to be valid", r)
		}
	}
}

func TestIsValidSupersessionReason_Unknown(t *testing.T) {
	if IsValidSupersessionReason("unknown-reason") {
		t.Error("expected unknown reason to be invalid")
	}
}

var baseSupersession = SecretSupersession{
	Mount:        "secret",
	Path:         "app/db-password",
	SupersededBy: "app/db-password-v2",
	Reason:       SupersessionReasonRotated,
	SupersededAt: time.Now(),
	Actor:        "ops-team",
}

func TestSecretSupersession_FullPath(t *testing.T) {
	s := baseSupersession
	if got := s.FullPath(); got != "secret/app/db-password" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretSupersession_Validate_Valid(t *testing.T) {
	if err := baseSupersession.Validate(); err != nil {
		t.Errorf("expected valid supersession, got: %v", err)
	}
}

func TestSecretSupersession_Validate_MissingMount(t *testing.T) {
	s := baseSupersession
	s.Mount = ""
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretSupersession_Validate_MissingPath(t *testing.T) {
	s := baseSupersession
	s.Path = ""
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretSupersession_Validate_MissingSupersededBy(t *testing.T) {
	s := baseSupersession
	s.SupersededBy = ""
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing superseded_by")
	}
}

func TestSecretSupersession_Validate_InvalidReason(t *testing.T) {
	s := baseSupersession
	s.Reason = "bad-reason"
	if err := s.Validate(); err == nil {
		t.Error("expected error for invalid reason")
	}
}

func TestSecretSupersession_Validate_MissingActor(t *testing.T) {
	s := baseSupersession
	s.Actor = ""
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing actor")
	}
}
