package vault

import (
	"testing"
)

var baseApproval = ApprovalRequest{
	ID:          "req-1",
	Mount:       "secret",
	Path:        "myapp/db",
	RequestedBy: "alice",
	Reason:      "rotate credentials",
}

func TestApprovalRequest_Validate_Valid(t *testing.T) {
	a := baseApproval
	if err := a.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestApprovalRequest_Validate_MissingMount(t *testing.T) {
	a := baseApproval
	a.Mount = ""
	if err := a.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestApprovalRequest_Validate_MissingPath(t *testing.T) {
	a := baseApproval
	a.Path = ""
	if err := a.Validate(); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestApprovalRequest_Validate_MissingRequestedBy(t *testing.T) {
	a := baseApproval
	a.RequestedBy = ""
	if err := a.Validate(); err == nil {
		t.Fatal("expected error for missing requested_by")
	}
}

func TestApprovalRequest_Validate_MissingReason(t *testing.T) {
	a := baseApproval
	a.Reason = ""
	if err := a.Validate(); err == nil {
		t.Fatal("expected error for missing reason")
	}
}

func TestApprovalRequest_FullPath(t *testing.T) {
	a := baseApproval
	want := "secret/myapp/db"
	if got := a.FullPath(); got != want {
		t.Fatalf("FullPath() = %q, want %q", got, want)
	}
}

func TestApprovalRequest_IsTerminal_Pending(t *testing.T) {
	a := baseApproval
	a.Status = ApprovalPending
	if a.IsTerminal() {
		t.Fatal("pending should not be terminal")
	}
}

func TestApprovalRequest_IsTerminal_Approved(t *testing.T) {
	a := baseApproval
	a.Status = ApprovalApproved
	if !a.IsTerminal() {
		t.Fatal("approved should be terminal")
	}
}
