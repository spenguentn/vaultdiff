package vault

import "testing"

func TestIsValidSensitivityLevel_Known(t *testing.T) {
	for _, lvl := range []SensitivityLevel{SensitivityLow, SensitivityMedium, SensitivityHigh, SensitivityCritical} {
		if !IsValidSensitivityLevel(lvl) {
			t.Errorf("expected %q to be valid", lvl)
		}
	}
}

func TestIsValidSensitivityLevel_Unknown(t *testing.T) {
	if IsValidSensitivityLevel("extreme") {
		t.Error("expected unknown level to be invalid")
	}
}

func TestSecretSensitivity_FullPath(t *testing.T) {
	s := SecretSensitivity{Mount: "secret", Path: "db/pass"}
	if got := s.FullPath(); got != "secret/db/pass" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretSensitivity_Validate_Valid(t *testing.T) {
	s := SecretSensitivity{Mount: "secret", Path: "db/pass", Level: SensitivityHigh, SetBy: "alice"}
	if err := s.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretSensitivity_Validate_MissingMount(t *testing.T) {
	s := SecretSensitivity{Path: "db/pass", Level: SensitivityLow, SetBy: "alice"}
	if err := s.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretSensitivity_Validate_InvalidLevel(t *testing.T) {
	s := SecretSensitivity{Mount: "secret", Path: "db/pass", Level: "extreme", SetBy: "alice"}
	if err := s.Validate(); err == nil {
		t.Error("expected error for invalid level")
	}
}

func TestSecretSensitivity_RequiresRedaction(t *testing.T) {
	for _, lvl := range []SensitivityLevel{SensitivityHigh, SensitivityCritical} {
		s := SecretSensitivity{Level: lvl}
		if !s.RequiresRedaction() {
			t.Errorf("expected %q to require redaction", lvl)
		}
	}
	for _, lvl := range []SensitivityLevel{SensitivityLow, SensitivityMedium} {
		s := SecretSensitivity{Level: lvl}
		if s.RequiresRedaction() {
			t.Errorf("expected %q not to require redaction", lvl)
		}
	}
}

func TestSensitivityRegistry_SetAndGet(t *testing.T) {
	r := NewSecretSensitivityRegistry()
	s := SecretSensitivity{Mount: "kv", Path: "app/token", Level: SensitivityCritical, SetBy: "bob"}
	if err := r.Set(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := r.Get("kv", "app/token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Level != SensitivityCritical {
		t.Errorf("expected critical, got %s", got.Level)
	}
}

func TestSensitivityRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretSensitivityRegistry()
	if _, err := r.Get("kv", "missing"); err == nil {
		t.Error("expected error for missing entry")
	}
}

func TestSensitivityRegistry_Remove(t *testing.T) {
	r := NewSecretSensitivityRegistry()
	s := SecretSensitivity{Mount: "kv", Path: "x", Level: SensitivityLow, SetBy: "eve"}
	_ = r.Set(s)
	r.Remove("kv", "x")
	if _, err := r.Get("kv", "x"); err == nil {
		t.Error("expected entry to be removed")
	}
}

func TestSensitivityRegistry_All(t *testing.T) {
	r := NewSecretSensitivityRegistry()
	_ = r.Set(SecretSensitivity{Mount: "kv", Path: "a", Level: SensitivityLow, SetBy: "u"})
	_ = r.Set(SecretSensitivity{Mount: "kv", Path: "b", Level: SensitivityMedium, SetBy: "u"})
	if len(r.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.All()))
	}
}
