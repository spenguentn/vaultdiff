package vault

import (
	"testing"
	"time"
)

var baseCustodian = SecretCustodian{
	Mount:      "secret",
	Path:       "app/db",
	Custodian:  "alice",
	Role:       CustodianRoleOwner,
	AssignedBy: "admin",
	AssignedAt: time.Now().UTC(),
}

func TestIsValidCustodianRole_Known(t *testing.T) {
	for _, r := range []CustodianRole{CustodianRoleOwner, CustodianRoleReviewer, CustodianRoleAuditor} {
		if !IsValidCustodianRole(r) {
			t.Errorf("expected %q to be valid", r)
		}
	}
}

func TestIsValidCustodianRole_Unknown(t *testing.T) {
	if IsValidCustodianRole("superuser") {
		t.Error("expected unknown role to be invalid")
	}
}

func TestSecretCustodian_FullPath(t *testing.T) {
	if got := baseCustodian.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretCustodian_Validate_Valid(t *testing.T) {
	if err := baseCustodian.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretCustodian_Validate_MissingMount(t *testing.T) {
	c := baseCustodian
	c.Mount = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretCustodian_Validate_InvalidRole(t *testing.T) {
	c := baseCustodian
	c.Role = "god"
	if err := c.Validate(); err == nil {
		t.Error("expected error for invalid role")
	}
}

func TestNewSecretCustodianRegistry_NotNil(t *testing.T) {
	if NewSecretCustodianRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestCustodianRegistry_Assign_And_Get(t *testing.T) {
	r := NewSecretCustodianRegistry()
	if err := r.Assign(baseCustodian); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get(baseCustodian.Mount, baseCustodian.Path, baseCustodian.Custodian)
	if !ok {
		t.Fatal("expected to find custodian")
	}
	if got.Role != CustodianRoleOwner {
		t.Errorf("unexpected role: %s", got.Role)
	}
}

func TestCustodianRegistry_Assign_SetsAssignedAt(t *testing.T) {
	r := NewSecretCustodianRegistry()
	c := baseCustodian
	c.AssignedAt = time.Time{}
	_ = r.Assign(c)
	got, _ := r.Get(c.Mount, c.Path, c.Custodian)
	if got.AssignedAt.IsZero() {
		t.Error("expected AssignedAt to be set")
	}
}

func TestCustodianRegistry_Remove(t *testing.T) {
	r := NewSecretCustodianRegistry()
	_ = r.Assign(baseCustodian)
	r.Remove(baseCustodian.Mount, baseCustodian.Path, baseCustodian.Custodian)
	if _, ok := r.Get(baseCustodian.Mount, baseCustodian.Path, baseCustodian.Custodian); ok {
		t.Error("expected custodian to be removed")
	}
}

func TestCustodianRegistry_List(t *testing.T) {
	r := NewSecretCustodianRegistry()
	_ = r.Assign(baseCustodian)
	second := baseCustodian
	second.Custodian = "bob"
	second.Role = CustodianRoleReviewer
	_ = r.Assign(second)
	list := r.List(baseCustodian.Mount, baseCustodian.Path)
	if len(list) != 2 {
		t.Errorf("expected 2 custodians, got %d", len(list))
	}
}
