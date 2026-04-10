package vault

import (
	"testing"
	"time"
)

func validLease(id string) LeaseInfo {
	return LeaseInfo{
		LeaseID:   id,
		Duration:  time.Hour,
		Renewable: true,
		IssuedAt:  time.Now(),
	}
}

func TestNewLeaseTracker_NotNil(t *testing.T) {
	tr := NewLeaseTracker()
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestLeaseTracker_Track_And_Get(t *testing.T) {
	tr := NewLeaseTracker()
	l := validLease("lease-001")
	if err := tr.Track(l); err != nil {
		t.Fatalf("unexpected track error: %v", err)
	}
	got, ok := tr.Get("lease-001")
	if !ok {
		t.Fatal("expected lease to be found")
	}
	if got.LeaseID != "lease-001" {
		t.Errorf("expected lease-001, got %q", got.LeaseID)
	}
}

func TestLeaseTracker_Track_InvalidLease(t *testing.T) {
	tr := NewLeaseTracker()
	bad := LeaseInfo{} // missing ID and duration
	if err := tr.Track(bad); err == nil {
		t.Error("expected validation error for invalid lease")
	}
}

func TestLeaseTracker_Get_NotFound(t *testing.T) {
	tr := NewLeaseTracker()
	_, ok := tr.Get("missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestLeaseTracker_Revoke_Success(t *testing.T) {
	tr := NewLeaseTracker()
	_ = tr.Track(validLease("lease-002"))
	if err := tr.Revoke("lease-002"); err != nil {
		t.Fatalf("unexpected revoke error: %v", err)
	}
	if tr.Count() != 0 {
		t.Error("expected empty tracker after revoke")
	}
}

func TestLeaseTracker_Revoke_NotFound(t *testing.T) {
	tr := NewLeaseTracker()
	if err := tr.Revoke("nonexistent"); err == nil {
		t.Error("expected error when revoking unknown lease")
	}
}

func TestLeaseTracker_Expired(t *testing.T) {
	tr := NewLeaseTracker()
	active := validLease("active")
	expired := LeaseInfo{
		LeaseID:  "expired",
		Duration: time.Millisecond,
		IssuedAt: time.Now().Add(-time.Second),
	}
	_ = tr.Track(active)
	_ = tr.Track(expired)

	exp := tr.Expired()
	if len(exp) != 1 {
		t.Fatalf("expected 1 expired lease, got %d", len(exp))
	}
	if exp[0].LeaseID != "expired" {
		t.Errorf("expected expired lease ID 'expired', got %q", exp[0].LeaseID)
	}
}

func TestLeaseTracker_Count(t *testing.T) {
	tr := NewLeaseTracker()
	_ = tr.Track(validLease("a"))
	_ = tr.Track(validLease("b"))
	if tr.Count() != 2 {
		t.Errorf("expected count 2, got %d", tr.Count())
	}
}
