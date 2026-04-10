package vault

import (
	"testing"
)

func basePolicy() Policy {
	return Policy{
		Name: "test-policy",
		Rules: []PolicyRule{
			{Path: "secret/data/app", Capabilities: []PolicyCapability{CapRead, CapList}},
			{Path: "secret/data/db", Capabilities: []PolicyCapability{CapRead}},
		},
	}
}

func TestPolicyRule_Validate_Valid(t *testing.T) {
	r := PolicyRule{Path: "secret/data/app", Capabilities: []PolicyCapability{CapRead}}
	if err := r.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPolicyRule_Validate_EmptyPath(t *testing.T) {
	r := PolicyRule{Path: "", Capabilities: []PolicyCapability{CapRead}}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestPolicyRule_Validate_NoCapabilities(t *testing.T) {
	r := PolicyRule{Path: "secret/data/app", Capabilities: nil}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for missing capabilities")
	}
}

func TestPolicyRule_HasCapability_True(t *testing.T) {
	r := PolicyRule{Path: "secret/data/app", Capabilities: []PolicyCapability{CapRead, CapList}}
	if !r.HasCapability(CapRead) {
		t.Fatal("expected HasCapability to return true for read")
	}
}

func TestPolicyRule_HasCapability_False(t *testing.T) {
	r := PolicyRule{Path: "secret/data/app", Capabilities: []PolicyCapability{CapRead}}
	if r.HasCapability(CapDelete) {
		t.Fatal("expected HasCapability to return false for delete")
	}
}

func TestPolicy_Validate_Valid(t *testing.T) {
	p := basePolicy()
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPolicy_Validate_EmptyName(t *testing.T) {
	p := basePolicy()
	p.Name = ""
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestPolicy_Validate_InvalidRule(t *testing.T) {
	p := basePolicy()
	p.Rules = append(p.Rules, PolicyRule{Path: "", Capabilities: []PolicyCapability{CapRead}})
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for invalid rule")
	}
}

func TestPolicy_RuleCount(t *testing.T) {
	p := basePolicy()
	if got := p.RuleCount(); got != 2 {
		t.Fatalf("expected 2 rules, got %d", got)
	}
}

func TestPolicy_AllowsRead_True(t *testing.T) {
	p := basePolicy()
	if !p.AllowsRead("secret/data/app") {
		t.Fatal("expected AllowsRead to return true")
	}
}

func TestPolicy_AllowsRead_False(t *testing.T) {
	p := basePolicy()
	if p.AllowsRead("secret/data/unknown") {
		t.Fatal("expected AllowsRead to return false for unknown path")
	}
}
