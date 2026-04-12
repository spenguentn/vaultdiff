package vault

import (
	"testing"
)

func baseCheckRequest() PolicyCheckRequest {
	return PolicyCheckRequest{
		Mount:      "secret",
		Path:       "myapp/db",
		Capability: "read",
	}
}

func TestPolicyCheckRequest_Validate_Valid(t *testing.T) {
	if err := baseCheckRequest().Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPolicyCheckRequest_Validate_MissingMount(t *testing.T) {
	req := baseCheckRequest()
	req.Mount = ""
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestPolicyCheckRequest_Validate_MissingPath(t *testing.T) {
	req := baseCheckRequest()
	req.Path = ""
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestPolicyCheckRequest_Validate_MissingCapability(t *testing.T) {
	req := baseCheckRequest()
	req.Capability = ""
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for missing capability")
	}
}

func TestPolicyCheckRequest_FullPath(t *testing.T) {
	req := baseCheckRequest()
	got := req.FullPath()
	want := "secret/myapp/db"
	if got != want {
		t.Errorf("FullPath() = %q, want %q", got, want)
	}
}

func TestCheckSecretPolicy_Allowed(t *testing.T) {
	rules := []PolicyRule{
		{Path: "secret/myapp/*", Capabilities: []string{"read", "list"}},
	}
	result, err := CheckSecretPolicy(baseCheckRequest(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsAllowed() {
		t.Errorf("expected allowed, got denied: %s", result.Reason)
	}
}

func TestCheckSecretPolicy_Denied_NoCapability(t *testing.T) {
	rules := []PolicyRule{
		{Path: "secret/myapp/*", Capabilities: []string{"list"}},
	}
	result, err := CheckSecretPolicy(baseCheckRequest(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsAllowed() {
		t.Error("expected denied, got allowed")
	}
}

func TestCheckSecretPolicy_NoMatchingRule(t *testing.T) {
	rules := []PolicyRule{
		{Path: "secret/other/*", Capabilities: []string{"read"}},
	}
	result, err := CheckSecretPolicy(baseCheckRequest(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsAllowed() {
		t.Error("expected denied for no matching rule")
	}
}

func TestCheckSecretPolicy_InvalidRequest(t *testing.T) {
	req := PolicyCheckRequest{}
	_, err := CheckSecretPolicy(req, nil)
	if err == nil {
		t.Fatal("expected error for invalid request")
	}
}

func TestCheckSecretPolicy_ExactPathMatch(t *testing.T) {
	rules := []PolicyRule{
		{Path: "secret/myapp/db", Capabilities: []string{"read", "write"}},
	}
	result, err := CheckSecretPolicy(baseCheckRequest(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsAllowed() {
		t.Errorf("expected allowed for exact match, reason: %s", result.Reason)
	}
}
