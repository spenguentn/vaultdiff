package vault

import (
	"errors"
	"testing"
	"time"
)

// stubWriter is a test double for LogicalWriter.
type stubWriter struct {
	resp *SecretResponse
	err  error
}

func (s *stubWriter) Write(_ string, _ map[string]interface{}) (*SecretResponse, error) {
	return s.resp, s.err
}

func TestRenewRequest_Validate_Valid(t *testing.T) {
	req := RenewRequest{LeaseID: "lease/abc123", Increment: 30 * time.Second}
	if err := req.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenewRequest_Validate_EmptyLeaseID(t *testing.T) {
	req := RenewRequest{}
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for empty lease ID")
	}
}

func TestNewRenewer_NilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil client")
		}
	}()
	NewRenewer(nil)
}

func TestNewRenewer_NotNil(t *testing.T) {
	w := &stubWriter{}
	r := NewRenewer(w)
	if r == nil {
		t.Fatal("expected non-nil Renewer")
	}
}

func TestRenewer_Renew_InvalidRequest(t *testing.T) {
	r := NewRenewer(&stubWriter{})
	res := r.Renew(RenewRequest{})
	if res.IsSuccess() {
		t.Fatal("expected failure for invalid request")
	}
}

func TestRenewer_Renew_WriterError(t *testing.T) {
	w := &stubWriter{err: errors.New("vault unavailable")}
	r := NewRenewer(w)
	res := r.Renew(RenewRequest{LeaseID: "lease/abc", Increment: 60 * time.Second})
	if res.IsSuccess() {
		t.Fatal("expected failure when writer returns error")
	}
	if res.Err == nil {
		t.Fatal("expected non-nil Err")
	}
}

func TestRenewer_Renew_Success(t *testing.T) {
	resp := &SecretResponse{
		LeaseID:       "lease/abc",
		LeaseDuration: 3600,
		Renewable:     true,
	}
	w := &stubWriter{resp: resp}
	r := NewRenewer(w)
	res := r.Renew(RenewRequest{LeaseID: "lease/abc", Increment: 60 * time.Second})
	if !res.IsSuccess() {
		t.Fatalf("expected success, got: %v", res.Err)
	}
	if res.LeaseID != "lease/abc" {
		t.Errorf("expected lease ID 'lease/abc', got %q", res.LeaseID)
	}
	if res.NewTTL != 3600*time.Second {
		t.Errorf("expected TTL 3600s, got %v", res.NewTTL)
	}
	if res.RenewedAt.IsZero() {
		t.Error("expected non-zero RenewedAt")
	}
}
