package vault

import (
	"testing"
	"time"
)

func baseCovenant() *SecretCovenant {
	return &SecretCovenant{
		Mount: "secret",
		Path:  "app/db",
		Type:  CovenantTypeShared,
		Owner: "team-platform",
	}
}

func TestIsValidCovenantType_Known(t *testing.T) {
	for _, tt := range []CovenantType{CovenantTypeShared, CovenantTypeExclusive, CovenantTypeReadOnly} {
		if !IsValidCovenantType(tt) {
			t.Errorf("expected %q to be valid", tt)
		}
	}
}

func TestIsValidCovenantType_Unknown(t *testing.T) {
	if IsValidCovenantType("bogus") {
		t.Error("expected 'bogus' to be invalid")
	}
}

func TestSecretCovenant_FullPath(t *testing.T) {
	c := baseCovenant()
	if got := c.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretCovenant_IsExpired_False(t *testing.T) {
	c := baseCovenant()
	future := time.Now().Add(time.Hour)
	c.ExpiresAt = &future
	if c.IsExpired() {
		t.Error("expected covenant not to be expired")
	}
}

func TestSecretCovenant_IsExpired_True(t *testing.T) {
	c := baseCovenant()
	past := time.Now().Add(-time.Hour)
	c.ExpiresAt = &past
	if !c.IsExpired() {
		t.Error("expected covenant to be expired")
	}
}

func TestSecretCovenant_IsExpired_NoExpiry(t *testing.T) {
	if baseCovenant().IsExpired() {
		t.Error("expected covenant with no expiry to not be expired")
	}
}

func TestSecretCovenant_Validate_Valid(t *testing.T) {
	if err := baseCovenant().Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretCovenant_Validate_MissingMount(t *testing.T) {
	c := baseCovenant()
	c.Mount = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretCovenant_Validate_MissingOwner(t *testing.T) {
	c := baseCovenant()
	c.Owner = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing owner")
	}
}

func TestSecretCovenant_Validate_InvalidType(t *testing.T) {
	c := baseCovenant()
	c.Type = "unknown"
	if err := c.Validate(); err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestCovenantRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretCovenantRegistry()
	if err := r.Set(baseCovenant()); err != nil {
		t.Fatalf("unexpected set error: %v", err)
	}
	got, err := r.Get("secret", "app/db")
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}
	if got.Owner != "team-platform" {
		t.Errorf("unexpected owner: %s", got.Owner)
	}
}

func TestCovenantRegistry_Set_SetsCreatedAt(t *testing.T) {
	r := NewSecretCovenantRegistry()
	c := baseCovenant()
	_ = r.Set(c)
	if c.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestCovenantRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretCovenantRegistry()
	if _, err := r.Get("secret", "missing"); err == nil {
		t.Error("expected error for missing covenant")
	}
}

func TestCovenantRegistry_Remove(t *testing.T) {
	r := NewSecretCovenantRegistry()
	_ = r.Set(baseCovenant())
	r.Remove("secret", "app/db")
	if _, err := r.Get("secret", "app/db"); err == nil {
		t.Error("expected error after removal")
	}
}

func TestCovenantRegistry_All(t *testing.T) {
	r := NewSecretCovenantRegistry()
	_ = r.Set(baseCovenant())
	if len(r.All()) != 1 {
		t.Errorf("expected 1 covenant, got %d", len(r.All()))
	}
}
