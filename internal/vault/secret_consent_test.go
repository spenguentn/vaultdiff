package vault

import (
	"testing"
	"time"
)

func baseConsent() SecretConsent {
	return SecretConsent{
		Mount:     "secret",
		Path:      "myapp/db",
		GrantedTo: "alice",
		GrantedBy: "admin",
		Status:    ConsentGranted,
		GrantedAt: time.Now().UTC(),
	}
}

func TestIsValidConsentStatus_Known(t *testing.T) {
	for _, s := range []ConsentStatus{ConsentPending, ConsentGranted, ConsentRevoked, ConsentExpired} {
		if !IsValidConsentStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidConsentStatus_Unknown(t *testing.T) {
	if IsValidConsentStatus("unknown") {
		t.Error("expected unknown status to be invalid")
	}
}

func TestSecretConsent_FullPath(t *testing.T) {
	c := baseConsent()
	if got := c.FullPath(); got != "secret/myapp/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretConsent_IsExpired_False(t *testing.T) {
	c := baseConsent()
	if c.IsExpired() {
		t.Error("expected consent without expiry to not be expired")
	}
}

func TestSecretConsent_IsExpired_True(t *testing.T) {
	c := baseConsent()
	past := time.Now().Add(-time.Hour)
	c.ExpiresAt = &past
	if !c.IsExpired() {
		t.Error("expected consent with past expiry to be expired")
	}
}

func TestSecretConsent_IsActive_True(t *testing.T) {
	c := baseConsent()
	if !c.IsActive() {
		t.Error("expected granted consent to be active")
	}
}

func TestSecretConsent_IsActive_False_Revoked(t *testing.T) {
	c := baseConsent()
	c.Status = ConsentRevoked
	if c.IsActive() {
		t.Error("expected revoked consent to be inactive")
	}
}

func TestSecretConsent_Validate_Valid(t *testing.T) {
	if err := baseConsent().Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretConsent_Validate_MissingMount(t *testing.T) {
	c := baseConsent()
	c.Mount = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretConsent_Validate_MissingGrantedTo(t *testing.T) {
	c := baseConsent()
	c.GrantedTo = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing granted_to")
	}
}

func TestSecretConsent_Validate_UnknownStatus(t *testing.T) {
	c := baseConsent()
	c.Status = "bad"
	if err := c.Validate(); err == nil {
		t.Error("expected error for unknown status")
	}
}
