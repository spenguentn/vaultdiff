package vault

import (
	"errors"
	"fmt"
	"time"
)

// CustodianRole represents the role a custodian holds over a secret.
type CustodianRole string

const (
	CustodianRoleOwner    CustodianRole = "owner"
	CustodianRoleReviewer CustodianRole = "reviewer"
	CustodianRoleAuditor  CustodianRole = "auditor"
)

// IsValidCustodianRole returns true if the role is a known custodian role.
func IsValidCustodianRole(r CustodianRole) bool {
	switch r {
	case CustodianRoleOwner, CustodianRoleReviewer, CustodianRoleAuditor:
		return true
	}
	return false
}

// SecretCustodian associates a named custodian with a specific secret path.
type SecretCustodian struct {
	Mount       string        `json:"mount"`
	Path        string        `json:"path"`
	Custodian   string        `json:"custodian"`
	Role        CustodianRole `json:"role"`
	AssignedAt  time.Time     `json:"assigned_at"`
	AssignedBy  string        `json:"assigned_by"`
}

// FullPath returns the canonical mount+path string.
func (c SecretCustodian) FullPath() string {
	return fmt.Sprintf("%s/%s", c.Mount, c.Path)
}

// Validate checks that the custodian record is well-formed.
func (c SecretCustodian) Validate() error {
	if c.Mount == "" {
		return errors.New("custodian: mount is required")
	}
	if c.Path == "" {
		return errors.New("custodian: path is required")
	}
	if c.Custodian == "" {
		return errors.New("custodian: custodian name is required")
	}
	if !IsValidCustodianRole(c.Role) {
		return fmt.Errorf("custodian: unknown role %q", c.Role)
	}
	if c.AssignedBy == "" {
		return errors.New("custodian: assigned_by is required")
	}
	return nil
}
