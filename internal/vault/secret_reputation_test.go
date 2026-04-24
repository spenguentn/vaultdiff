package vault

import (
	"testing"
	"time"
)

func TestIsValidReputationLevel_Known(t *testing.T) {
	for _, lvl := range []ReputationLevel{
		ReputationExemplary, ReputationGood, ReputationNeutral,
		ReputationSuspect, ReputationCompromised,
	} {
		if !IsValidReputationLevel(lvl) {
			t.Errorf("expected %q to be valid", lvl)
		}
	}
}

func TestIsValidReputationLevel_Unknown(t *testing.T) {
	if IsValidReputationLevel("legendary") {
		t.Error("expected 'legendary' to be invalid")
	}
}

func TestSecretReputation_FullPath(t *testing.T) {
	r := SecretReputation{Mount: "secret", Path: "app/db"}
	if got := r.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretReputation_Validate_Valid(t *testing.T) {
	r := SecretReputation{
		Mount:      "secret",
		Path:       "app/key",
		Level:      ReputationGood,
		AssessedBy: "alice",
	}
	if err := r.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretReputation_Validate_MissingMount(t *testing.T) {
	r := SecretReputation{Path: "app/key", Level: ReputationGood, AssessedBy: "alice"}
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretReputation_Validate_InvalidLevel(t *testing.T) {
	r := SecretReputation{Mount: "secret", Path: "app/key", Level: "unknown", AssessedBy: "alice"}
	if err := r.Validate(); err == nil {
		t.Error("expected error for invalid level")
	}
}

func TestNewSecretReputationRegistry_NotNil(t *testing.T) {
	if NewSecretReputationRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestReputationRegistry_Set_And_Get(t *testing.T) {
	reg := NewSecretReputationRegistry()
	r := SecretReputation{
		Mount: "secret", Path: "svc/token",
		Level: ReputationNeutral, AssessedBy: "bot",
	}
	if err := reg.Set(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("secret", "svc/token")
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.Level != ReputationNeutral {
		t.Errorf("unexpected level: %s", got.Level)
	}
}

func TestReputationRegistry_Set_SetsAssessedAt(t *testing.T) {
	reg := NewSecretReputationRegistry()
	r := SecretReputation{
		Mount: "secret", Path: "svc/key",
		Level: ReputationGood, AssessedBy: "ops",
	}
	_ = reg.Set(r)
	got, _ := reg.Get("secret", "svc/key")
	if got.AssessedAt.IsZero() {
		t.Error("expected AssessedAt to be stamped")
	}
}

func TestReputationRegistry_Set_Invalid(t *testing.T) {
	reg := NewSecretReputationRegistry()
	if err := reg.Set(SecretReputation{}); err == nil {
		t.Error("expected validation error")
	}
}

func TestReputationRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretReputationRegistry()
	_, ok := reg.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestReputationRegistry_Remove(t *testing.T) {
	reg := NewSecretReputationRegistry()
	r := SecretReputation{
		Mount: "secret", Path: "app/x",
		Level: ReputationExemplary, AssessedBy: "admin",
		AssessedAt: time.Now(),
	}
	_ = reg.Set(r)
	reg.Remove("secret", "app/x")
	if _, ok := reg.Get("secret", "app/x"); ok {
		t.Error("expected record to be removed")
	}
}
