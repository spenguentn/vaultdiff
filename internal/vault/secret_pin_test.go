package vault

import (
	"testing"
	"time"
)

var basePin = SecretPin{
	Mount:    "secret",
	Path:     "app/db",
	Version:  3,
	PinnedBy: "alice",
	Reason:   "stable release freeze",
}

func TestSecretPin_FullPath(t *testing.T) {
	got := basePin.FullPath()
	want := "secret/app/db"
	if got != want {
		t.Errorf("FullPath() = %q, want %q", got, want)
	}
}

func TestSecretPin_IsExpired_False(t *testing.T) {
	p := basePin
	p.ExpiresAt = time.Now().Add(time.Hour)
	if p.IsExpired() {
		t.Error("expected pin not to be expired")
	}
}

func TestSecretPin_IsExpired_True(t *testing.T) {
	p := basePin
	p.ExpiresAt = time.Now().Add(-time.Minute)
	if !p.IsExpired() {
		t.Error("expected pin to be expired")
	}
}

func TestSecretPin_IsExpired_NoExpiry(t *testing.T) {
	if basePin.IsExpired() {
		t.Error("pin with zero ExpiresAt should never be expired")
	}
}

func TestSecretPin_Validate_Valid(t *testing.T) {
	if err := basePin.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretPin_Validate_MissingMount(t *testing.T) {
	p := basePin
	p.Mount = ""
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretPin_Validate_MissingPath(t *testing.T) {
	p := basePin
	p.Path = ""
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretPin_Validate_ZeroVersion(t *testing.T) {
	p := basePin
	p.Version = 0
	if err := p.Validate(); err == nil {
		t.Error("expected error for zero version")
	}
}

func TestSecretPin_Validate_MissingPinnedBy(t *testing.T) {
	p := basePin
	p.PinnedBy = ""
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing pinned_by")
	}
}

func TestSecretPin_Validate_MissingReason(t *testing.T) {
	p := basePin
	p.Reason = ""
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing reason")
	}
}
