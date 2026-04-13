package vault

import (
	"testing"
)

func TestIsValidTrustLevel_Known(t *testing.T) {
	levels := []TrustLevel{TrustLevelUntrusted, TrustLevelLow, TrustLevelMedium, TrustLevelHigh, TrustLevelVerified}
	for _, l := range levels {
		if !IsValidTrustLevel(l) {
			t.Errorf("expected %q to be valid", l)
		}
	}
}

func TestIsValidTrustLevel_Unknown(t *testing.T) {
	if IsValidTrustLevel("supreme") {
		t.Error("expected 'supreme' to be invalid")
	}
}

func TestSecretTrust_FullPath(t *testing.T) {
	tr := &SecretTrust{Mount: "secret", Path: "app/db"}
	if got := tr.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretTrust_Validate_Valid(t *testing.T) {
	tr := &SecretTrust{
		Mount:      "secret",
		Path:       "app/db",
		Level:      TrustLevelHigh,
		AssignedBy: "alice",
	}
	if err := tr.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretTrust_Validate_MissingMount(t *testing.T) {
	tr := &SecretTrust{Path: "app/db", Level: TrustLevelLow, AssignedBy: "alice"}
	if err := tr.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretTrust_Validate_MissingPath(t *testing.T) {
	tr := &SecretTrust{Mount: "secret", Level: TrustLevelLow, AssignedBy: "alice"}
	if err := tr.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretTrust_Validate_MissingAssignedBy(t *testing.T) {
	tr := &SecretTrust{Mount: "secret", Path: "app/db", Level: TrustLevelMedium}
	if err := tr.Validate(); err == nil {
		t.Error("expected error for missing assigned_by")
	}
}

func TestSecretTrust_Validate_InvalidLevel(t *testing.T) {
	tr := &SecretTrust{Mount: "secret", Path: "app/db", Level: "godlike", AssignedBy: "alice"}
	if err := tr.Validate(); err == nil {
		t.Error("expected error for invalid trust level")
	}
}
