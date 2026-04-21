package vault

import (
	"testing"
	"time"
)

func TestIsValidRevocationReason_Known(t *testing.T) {
	for _, r := range []string{"compromised", "expired", "rotated", "decommissioned", "policy_violation"} {
		if !IsValidRevocationReason(r) {
			t.Errorf("expected %q to be valid", r)
		}
	}
}

func TestIsValidRevocationReason_Unknown(t *testing.T) {
	if IsValidRevocationReason("unknown_reason") {
		t.Error("expected unknown_reason to be invalid")
	}
}

func TestRevocationRecord_FullPath(t *testing.T) {
	rec := RevocationRecord{Mount: "secret", Path: "app/key"}
	if got := rec.FullPath(); got != "secret/app/key" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestRevocationRecord_IsExpired_False(t *testing.T) {
	future := time.Now().Add(time.Hour)
	rec := RevocationRecord{ExpiresAt: &future}
	if rec.IsExpired() {
		t.Error("expected not expired")
	}
}

func TestRevocationRecord_IsExpired_True(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	rec := RevocationRecord{ExpiresAt: &past}
	if !rec.IsExpired() {
		t.Error("expected expired")
	}
}

func TestRevocationRecord_IsExpired_NoExpiry(t *testing.T) {
	rec := RevocationRecord{}
	if rec.IsExpired() {
		t.Error("expected not expired when no expiry set")
	}
}

func TestRevocationRecord_Validate_Valid(t *testing.T) {
	rec := RevocationRecord{
		Mount: "secret", Path: "app/key",
		RevokedBy: "admin", Reason: "rotated",
		RevokedAt: time.Now(),
	}
	if err := rec.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRevocationRecord_Validate_MissingMount(t *testing.T) {
	rec := RevocationRecord{Path: "app/key", RevokedBy: "admin", Reason: "rotated", RevokedAt: time.Now()}
	if err := rec.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestRevocationRecord_Validate_UnknownReason(t *testing.T) {
	rec := RevocationRecord{
		Mount: "secret", Path: "app/key",
		RevokedBy: "admin", Reason: "bad_reason",
		RevokedAt: time.Now(),
	}
	if err := rec.Validate(); err == nil {
		t.Error("expected error for unknown reason")
	}
}
