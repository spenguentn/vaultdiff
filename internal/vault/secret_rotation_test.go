package vault

import (
	"testing"
	"time"
)

func basePolicy() RotationPolicy {
	return RotationPolicy{
		Mount:    "secret",
		Path:     "myapp/db",
		Interval: 24 * time.Hour,
		Enabled:  true,
	}
}

func TestRotationPolicy_Validate_Valid(t *testing.T) {
	if err := basePolicy().Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRotationPolicy_Validate_MissingMount(t *testing.T) {
	p := basePolicy()
	p.Mount = ""
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestRotationPolicy_Validate_MissingPath(t *testing.T) {
	p := basePolicy()
	p.Path = ""
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestRotationPolicy_Validate_ZeroInterval(t *testing.T) {
	p := basePolicy()
	p.Interval = 0
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestRotationPolicy_IsDue_NotEnabled(t *testing.T) {
	p := basePolicy()
	p.Enabled = false
	if p.IsDue() {
		t.Fatal("expected IsDue to be false when disabled")
	}
}

func TestRotationPolicy_IsDue_ZeroLastRotated(t *testing.T) {
	p := basePolicy()
	if !p.IsDue() {
		t.Fatal("expected IsDue to be true when never rotated")
	}
}

func TestRotationPolicy_IsDue_RecentRotation(t *testing.T) {
	p := basePolicy()
	p.LastRotatedAt = time.Now()
	if p.IsDue() {
		t.Fatal("expected IsDue to be false after recent rotation")
	}
}

func TestRotationResult_IsSuccess_True(t *testing.T) {
	r := RotationResult{Mount: "secret", Path: "myapp/db", RotatedAt: time.Now()}
	if !r.IsSuccess() {
		t.Fatal("expected IsSuccess to be true")
	}
}

func TestRotationResult_IsSuccess_WithErr(t *testing.T) {
	r := RotationResult{Err: errTest}
	if r.IsSuccess() {
		t.Fatal("expected IsSuccess to be false when Err is set")
	}
}

func TestNewRotationScheduler_NotNil(t *testing.T) {
	if s := NewRotationScheduler(); s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestRotationScheduler_Register_And_Count(t *testing.T) {
	s := NewRotationScheduler()
	if err := s.Register(basePolicy()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Count() != 1 {
		t.Fatalf("expected count 1, got %d", s.Count())
	}
}

func TestRotationScheduler_Register_Invalid(t *testing.T) {
	s := NewRotationScheduler()
	p := basePolicy()
	p.Mount = ""
	if err := s.Register(p); err == nil {
		t.Fatal("expected error for invalid policy")
	}
}

func TestRotationScheduler_Remove_Found(t *testing.T) {
	s := NewRotationScheduler()
	_ = s.Register(basePolicy())
	if err := s.Remove("secret", "myapp/db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Count() != 0 {
		t.Fatal("expected count 0 after removal")
	}
}

func TestRotationScheduler_Remove_NotFound(t *testing.T) {
	s := NewRotationScheduler()
	if err := s.Remove("secret", "missing"); err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestRotationScheduler_DuePolicies(t *testing.T) {
	s := NewRotationScheduler()
	p := basePolicy()
	_ = s.Register(p)
	due := s.DuePolicies()
	if len(due) != 1 {
		t.Fatalf("expected 1 due policy, got %d", len(due))
	}
}
