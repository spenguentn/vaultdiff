package vault

import (
	"testing"
	"time"
)

func TestParseLease_Valid(t *testing.T) {
	l, err := ParseLease("lease-abc-123", 3600, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.LeaseID != "lease-abc-123" {
		t.Errorf("expected lease ID 'lease-abc-123', got %q", l.LeaseID)
	}
	if l.Duration != 3600*time.Second {
		t.Errorf("expected duration 3600s, got %v", l.Duration)
	}
	if !l.Renewable {
		t.Error("expected renewable to be true")
	}
}

func TestParseLease_EmptyID(t *testing.T) {
	_, err := ParseLease("", 3600, true)
	if err == nil {
		t.Fatal("expected error for empty lease ID")
	}
}

func TestParseLease_ZeroDuration(t *testing.T) {
	_, err := ParseLease("lease-xyz", 0, false)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestLeaseInfo_ExpiresAt(t *testing.T) {
	now := time.Now()
	l := LeaseInfo{LeaseID: "id", Duration: 60 * time.Second, IssuedAt: now}
	expected := now.Add(60 * time.Second)
	if !l.ExpiresAt().Equal(expected) {
		t.Errorf("ExpiresAt mismatch: got %v, want %v", l.ExpiresAt(), expected)
	}
}

func TestLeaseInfo_IsExpired_False(t *testing.T) {
	l := LeaseInfo{LeaseID: "id", Duration: time.Hour, IssuedAt: time.Now()}
	if l.IsExpired() {
		t.Error("expected lease to not be expired")
	}
}

func TestLeaseInfo_IsExpired_True(t *testing.T) {
	l := LeaseInfo{
		LeaseID:  "id",
		Duration: time.Second,
		IssuedAt: time.Now().Add(-2 * time.Second),
	}
	if !l.IsExpired() {
		t.Error("expected lease to be expired")
	}
}

func TestLeaseInfo_TTLRemaining_Positive(t *testing.T) {
	l := LeaseInfo{LeaseID: "id", Duration: time.Hour, IssuedAt: time.Now()}
	if l.TTLRemaining() <= 0 {
		t.Error("expected positive TTL remaining")
	}
}

func TestLeaseInfo_TTLRemaining_Zero_WhenExpired(t *testing.T) {
	l := LeaseInfo{
		LeaseID:  "id",
		Duration: time.Millisecond,
		IssuedAt: time.Now().Add(-time.Second),
	}
	if l.TTLRemaining() != 0 {
		t.Errorf("expected 0 TTL remaining for expired lease, got %v", l.TTLRemaining())
	}
}

func TestLeaseInfo_Validate_Valid(t *testing.T) {
	l := LeaseInfo{LeaseID: "abc", Duration: time.Minute}
	if err := l.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestLeaseInfo_Validate_MissingID(t *testing.T) {
	l := LeaseInfo{Duration: time.Minute}
	if err := l.Validate(); err == nil {
		t.Error("expected error for missing lease ID")
	}
}

func TestLeaseInfo_Validate_ZeroDuration(t *testing.T) {
	l := LeaseInfo{LeaseID: "abc"}
	if err := l.Validate(); err == nil {
		t.Error("expected error for zero duration")
	}
}
