package vault

import (
	"strings"
	"testing"
)

func validReq() RollbackRequest {
	return RollbackRequest{Mount: "secret", Path: "app/db", Version: 2}
}

func TestNewRollbackPlan_Empty(t *testing.T) {
	p := NewRollbackPlan()
	if !p.IsEmpty() {
		t.Fatal("new plan should be empty")
	}
	if p.Len() != 0 {
		t.Fatalf("expected len 0, got %d", p.Len())
	}
}

func TestRollbackPlan_Add_Valid(t *testing.T) {
	p := NewRollbackPlan()
	if err := p.Add(validReq()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 1 {
		t.Fatalf("expected len 1, got %d", p.Len())
	}
}

func TestRollbackPlan_Add_MissingMount(t *testing.T) {
	p := NewRollbackPlan()
	err := p.Add(RollbackRequest{Path: "app/db", Version: 1})
	if err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestRollbackPlan_Add_MissingPath(t *testing.T) {
	p := NewRollbackPlan()
	err := p.Add(RollbackRequest{Mount: "secret", Version: 1})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestRollbackPlan_Add_ZeroVersion(t *testing.T) {
	p := NewRollbackPlan()
	err := p.Add(RollbackRequest{Mount: "secret", Path: "app/db", Version: 0})
	if err == nil {
		t.Fatal("expected error for version 0")
	}
}

func TestRollbackPlan_Requests_IsCopy(t *testing.T) {
	p := NewRollbackPlan()
	_ = p.Add(validReq())
	reqs := p.Requests()
	reqs[0].Version = 999
	if p.Requests()[0].Version == 999 {
		t.Fatal("Requests() should return a copy, not a reference")
	}
}

func TestRollbackPlan_Describe_Empty(t *testing.T) {
	p := NewRollbackPlan()
	if !strings.Contains(p.Describe(), "empty") {
		t.Fatal("expected 'empty' in description of empty plan")
	}
}

func TestRollbackPlan_Describe_WithRequests(t *testing.T) {
	p := NewRollbackPlan()
	_ = p.Add(validReq())
	desc := p.Describe()
	if !strings.Contains(desc, "1 operation") {
		t.Fatalf("unexpected describe output: %s", desc)
	}
	if !strings.Contains(desc, "secret/app/db") {
		t.Fatalf("expected path in describe output: %s", desc)
	}
}
