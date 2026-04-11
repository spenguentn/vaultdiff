package vault

import (
	"testing"
)

func newTestApproval() *ApprovalRequest {
	a := baseApproval
	return &a
}

func TestNewSecretApprovalRegistry_NotNil(t *testing.T) {
	r := NewSecretApprovalRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestApprovalRegistry_Submit_And_Get(t *testing.T) {
	r := NewSecretApprovalRegistry()
	req := newTestApproval()
	if err := r.Submit(req); err != nil {
		t.Fatalf("Submit() error: %v", err)
	}
	got, ok := r.Get(req.Mount, req.Path)
	if !ok {
		t.Fatal("expected to find submitted request")
	}
	if got.Status != ApprovalPending {
		t.Fatalf("expected pending, got %s", got.Status)
	}
}

func TestApprovalRegistry_Submit_Invalid(t *testing.T) {
	r := NewSecretApprovalRegistry()
	req := newTestApproval()
	req.Mount = ""
	if err := r.Submit(req); err == nil {
		t.Fatal("expected error for invalid request")
	}
}

func TestApprovalRegistry_Submit_DuplicatePending(t *testing.T) {
	r := NewSecretApprovalRegistry()
	r.Submit(newTestApproval())
	if err := r.Submit(newTestApproval()); err == nil {
		t.Fatal("expected error for duplicate pending request")
	}
}

func TestApprovalRegistry_Review_Approved(t *testing.T) {
	r := NewSecretApprovalRegistry()
	req := newTestApproval()
	r.Submit(req)
	if err := r.Review(req.Mount, req.Path, "bob", true); err != nil {
		t.Fatalf("Review() error: %v", err)
	}
	got, _ := r.Get(req.Mount, req.Path)
	if got.Status != ApprovalApproved {
		t.Fatalf("expected approved, got %s", got.Status)
	}
	if got.ReviewedBy != "bob" {
		t.Fatalf("expected reviewer bob, got %s", got.ReviewedBy)
	}
}

func TestApprovalRegistry_Review_Rejected(t *testing.T) {
	r := NewSecretApprovalRegistry()
	req := newTestApproval()
	r.Submit(req)
	r.Review(req.Mount, req.Path, "carol", false)
	got, _ := r.Get(req.Mount, req.Path)
	if got.Status != ApprovalRejected {
		t.Fatalf("expected rejected, got %s", got.Status)
	}
}

func TestApprovalRegistry_Review_NotFound(t *testing.T) {
	r := NewSecretApprovalRegistry()
	if err := r.Review("secret", "missing", "bob", true); err == nil {
		t.Fatal("expected error for missing request")
	}
}

func TestApprovalRegistry_Revoke(t *testing.T) {
	r := NewSecretApprovalRegistry()
	req := newTestApproval()
	r.Submit(req)
	r.Review(req.Mount, req.Path, "bob", true)
	if err := r.Revoke(req.Mount, req.Path); err != nil {
		t.Fatalf("Revoke() error: %v", err)
	}
	got, _ := r.Get(req.Mount, req.Path)
	if got.Status != ApprovalRevoked {
		t.Fatalf("expected revoked, got %s", got.Status)
	}
}

func TestApprovalRegistry_Revoke_NotApproved(t *testing.T) {
	r := NewSecretApprovalRegistry()
	req := newTestApproval()
	r.Submit(req)
	if err := r.Revoke(req.Mount, req.Path); err == nil {
		t.Fatal("expected error revoking non-approved request")
	}
}

func TestApprovalRegistry_All(t *testing.T) {
	r := NewSecretApprovalRegistry()
	r.Submit(newTestApproval())
	if len(r.All()) != 1 {
		t.Fatalf("expected 1 request, got %d", len(r.All()))
	}
}
