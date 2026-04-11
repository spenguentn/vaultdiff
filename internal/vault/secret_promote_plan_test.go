package vault

import (
	"testing"
)

func TestNewPromotePlan_Empty(t *testing.T) {
	p := NewPromotePlan()
	if p == nil {
		t.Fatal("expected non-nil plan")
	}
	if p.Len() != 0 {
		t.Fatalf("expected 0 requests, got %d", p.Len())
	}
}

func TestPromotePlan_Add_Valid(t *testing.T) {
	p := NewPromotePlan()
	if err := p.Add(validPromote); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 1 {
		t.Fatalf("expected 1, got %d", p.Len())
	}
}

func TestPromotePlan_Add_Invalid(t *testing.T) {
	p := NewPromotePlan()
	bad := PromoteRequest{} // missing required fields
	if err := p.Add(bad); err == nil {
		t.Fatal("expected error for invalid request")
	}
	if p.Len() != 0 {
		t.Fatal("plan should remain empty after failed add")
	}
}

func TestPromotePlan_Requests_ReturnsCopy(t *testing.T) {
	p := NewPromotePlan()
	_ = p.Add(validPromote)
	reqs := p.Requests()
	reqs[0].SourceMount = "mutated"
	if p.Requests()[0].SourceMount == "mutated" {
		t.Fatal("Requests should return a copy, not a reference")
	}
}

func TestPromotePlan_Validate_Empty(t *testing.T) {
	p := NewPromotePlan()
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty plan")
	}
}

func TestPromotePlan_Validate_NonEmpty(t *testing.T) {
	p := NewPromotePlan()
	_ = p.Add(validPromote)
	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
