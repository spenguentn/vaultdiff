package vault

import (
	"testing"
	"time"
)

var baseBadge = SecretBadge{
	Mount:     "secret",
	Path:      "apps/api/key",
	Badge:     BadgeTypeVerified,
	AwardedBy: "admin",
	AwardedAt: time.Now(),
}

func TestIsValidBadgeType_Known(t *testing.T) {
	for _, bt := range []BadgeType{BadgeTypeCompliant, BadgeTypeRotated, BadgeTypeVerified, BadgeTypeDeprecated} {
		if !IsValidBadgeType(bt) {
			t.Errorf("expected %q to be valid", bt)
		}
	}
}

func TestIsValidBadgeType_Unknown(t *testing.T) {
	if IsValidBadgeType(BadgeType("unknown")) {
		t.Error("expected unknown badge type to be invalid")
	}
}

func TestSecretBadge_FullPath(t *testing.T) {
	got := baseBadge.FullPath()
	want := "secret/apps/api/key"
	if got != want {
		t.Errorf("FullPath() = %q, want %q", got, want)
	}
}

func TestSecretBadge_Validate_Valid(t *testing.T) {
	if err := baseBadge.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretBadge_Validate_MissingMount(t *testing.T) {
	b := baseBadge
	b.Mount = ""
	if err := b.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretBadge_Validate_MissingPath(t *testing.T) {
	b := baseBadge
	b.Path = ""
	if err := b.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretBadge_Validate_InvalidBadgeType(t *testing.T) {
	b := baseBadge
	b.Badge = BadgeType("gold-star")
	if err := b.Validate(); err == nil {
		t.Error("expected error for invalid badge type")
	}
}

func TestSecretBadge_Validate_MissingAwardedBy(t *testing.T) {
	b := baseBadge
	b.AwardedBy = ""
	if err := b.Validate(); err == nil {
		t.Error("expected error for missing awarded_by")
	}
}

func TestSecretBadge_Validate_ZeroAwardedAt(t *testing.T) {
	b := baseBadge
	b.AwardedAt = time.Time{}
	if err := b.Validate(); err == nil {
		t.Error("expected error for zero awarded_at")
	}
}
