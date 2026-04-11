package vault

import (
	"testing"
	"time"
)

var baseExpiry = ExpiryPolicy{
	Mount:      "secret",
	Path:       "myapp/api-key",
	ExpiresAt:  time.Now().Add(72 * time.Hour),
	WarnBefore: 24 * time.Hour,
	Owner:      "team-platform",
}

func TestExpiryPolicy_Validate_Valid(t *testing.T) {
	if err := baseExpiry.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExpiryPolicy_Validate_MissingMount(t *testing.T) {
	p := baseExpiry
	p.Mount = ""
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestExpiryPolicy_Validate_MissingPath(t *testing.T) {
	p := baseExpiry
	p.Path = ""
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestExpiryPolicy_Validate_ZeroExpiresAt(t *testing.T) {
	p := baseExpiry
	p.ExpiresAt = time.Time{}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero expires_at")
	}
}

func TestExpiryPolicy_FullPath(t *testing.T) {
	got := baseExpiry.FullPath()
	want := "secret/myapp/api-key"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestExpiryPolicy_IsExpired_False(t *testing.T) {
	if baseExpiry.IsExpired(time.Now()) {
		t.Fatal("expected not expired")
	}
}

func TestExpiryPolicy_IsExpired_True(t *testing.T) {
	p := baseExpiry
	p.ExpiresAt = time.Now().Add(-1 * time.Hour)
	if !p.IsExpired(time.Now()) {
		t.Fatal("expected expired")
	}
}

func TestExpiryPolicy_IsExpiringSoon_True(t *testing.T) {
	p := baseExpiry
	p.ExpiresAt = time.Now().Add(12 * time.Hour)
	if !p.IsExpiringSoon(time.Now()) {
		t.Fatal("expected expiring soon")
	}
}

func TestExpiryStatus_String_OK(t *testing.T) {
	s := ExpiryStatus{Policy: baseExpiry}
	if s.String() != "ok" {
		t.Fatalf("got %q, want \"ok\"", s.String())
	}
}

func TestExpiryStatus_String_Expired(t *testing.T) {
	s := ExpiryStatus{Expired: true}
	if s.String() != "expired" {
		t.Fatalf("got %q, want \"expired\"", s.String())
	}
}

func TestExpiryStatus_String_Soon(t *testing.T) {
	s := ExpiryStatus{Soon: true}
	if s.String() != "expiring-soon" {
		t.Fatalf("got %q, want \"expiring-soon\"", s.String())
	}
}
