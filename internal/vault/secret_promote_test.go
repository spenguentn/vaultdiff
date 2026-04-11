package vault

import (
	"testing"
)

var validPromote = PromoteRequest{
	SourceMount: "secret",
	SourcePath:  "app/config",
	DestMount:   "secret-prod",
	DestPath:    "app/config",
}

func TestPromoteRequest_Validate_Valid(t *testing.T) {
	if err := validPromote.Validate(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestPromoteRequest_Validate_MissingSourceMount(t *testing.T) {
	r := validPromote
	r.SourceMount = ""
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for missing source mount")
	}
}

func TestPromoteRequest_Validate_MissingDestPath(t *testing.T) {
	r := validPromote
	r.DestPath = ""
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for missing dest path")
	}
}

func TestPromoteRequest_Validate_NegativeVersion(t *testing.T) {
	r := validPromote
	r.Version = -1
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for negative version")
	}
}

func TestPromoteResult_IsSuccess_True(t *testing.T) {
	res := PromoteResult{Request: validPromote}
	if !res.IsSuccess() {
		t.Fatal("expected success")
	}
}

func TestPromoteResult_IsSuccess_WithErr(t *testing.T) {
	res := PromoteResult{Request: validPromote, Err: errSentinel("boom")}
	if res.IsSuccess() {
		t.Fatal("expected failure")
	}
}

func TestPromoteResult_String_DryRun(t *testing.T) {
	r := validPromote
	r.DryRun = true
	res := PromoteResult{Request: r}
	s := res.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
	if len(s) < 8 {
		t.Fatalf("unexpected string: %s", s)
	}
}

func TestPromoteResult_String_Err(t *testing.T) {
	res := PromoteResult{Request: validPromote, Err: errSentinel("oops")}
	if res.String() == "" {
		t.Fatal("expected non-empty error string")
	}
}
