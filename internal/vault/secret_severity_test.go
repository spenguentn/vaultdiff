package vault

import "testing"

func TestIsValidSeverityLevel_Known(t *testing.T) {
	for _, lvl := range []SeverityLevel{SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical} {
		if !IsValidSeverityLevel(lvl) {
			t.Errorf("expected %q to be valid", lvl)
		}
	}
}

func TestIsValidSeverityLevel_Unknown(t *testing.T) {
	if IsValidSeverityLevel("extreme") {
		t.Error("expected 'extreme' to be invalid")
	}
}

func TestSecretSeverity_FullPath(t *testing.T) {
	s := SecretSeverity{Mount: "secret", Path: "app/db"}
	if got := s.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretSeverity_Validate_Valid(t *testing.T) {
	s := SecretSeverity{Mount: "secret", Path: "app/db", Level: SeverityHigh, AssignedBy: "alice"}
	if err := s.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretSeverity_Validate_MissingMount(t *testing.T) {
	s := SecretSeverity{Path: "app/db", Level: SeverityLow, AssignedBy: "alice"}
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretSeverity_Validate_MissingPath(t *testing.T) {
	s := SecretSeverity{Mount: "secret", Level: SeverityLow, AssignedBy: "alice"}
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretSeverity_Validate_InvalidLevel(t *testing.T) {
	s := SecretSeverity{Mount: "secret", Path: "app/db", Level: "extreme", AssignedBy: "alice"}
	if err := s.Validate(); err == nil {
		t.Error("expected error for invalid level")
	}
}

func TestSecretSeverity_Validate_MissingAssignedBy(t *testing.T) {
	s := SecretSeverity{Mount: "secret", Path: "app/db", Level: SeverityCritical}
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing assigned_by")
	}
}
