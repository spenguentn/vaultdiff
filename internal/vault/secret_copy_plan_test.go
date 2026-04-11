package vault

import "testing"

func TestNewCopyPlan_Empty(t *testing.T) {
	p := NewCopyPlan()
	if p.Len() != 0 {
		t.Fatalf("expected 0 requests, got %d", p.Len())
	}
}

func TestCopyPlan_Add_Valid(t *testing.T) {
	p := NewCopyPlan()
	err := p.Add(CopyRequest{SourceMount: "kv", SourcePath: "a", DestMount: "kv", DestPath: "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 1 {
		t.Fatalf("expected 1 request, got %d", p.Len())
	}
}

func TestCopyPlan_Add_Invalid(t *testing.T) {
	p := NewCopyPlan()
	err := p.Add(CopyRequest{SourcePath: "a", DestMount: "kv", DestPath: "b"})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if p.Len() != 0 {
		t.Fatal("invalid request should not be added")
	}
}

func TestCopyPlan_Requests_ReturnsCopy(t *testing.T) {
	p := NewCopyPlan()
	_ = p.Add(CopyRequest{SourceMount: "kv", SourcePath: "a", DestMount: "kv", DestPath: "b"})
	reqs := p.Requests()
	if len(reqs) != 1 {
		t.Fatalf("expected 1, got %d", len(reqs))
	}
	// Mutating the returned slice must not affect the plan.
	reqs[0].SourcePath = "mutated"
	if p.Requests()[0].SourcePath == "mutated" {
		t.Fatal("Requests() should return an independent copy")
	}
}

func TestCopyPlan_Validate_Empty(t *testing.T) {
	p := NewCopyPlan()
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty plan")
	}
}

func TestCopyPlan_Validate_NonEmpty(t *testing.T) {
	p := NewCopyPlan()
	_ = p.Add(CopyRequest{SourceMount: "kv", SourcePath: "a", DestMount: "kv", DestPath: "b"})
	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
