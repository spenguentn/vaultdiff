package vault

import (
	"errors"
	"fmt"
	"time"
)

// DelegationScope defines the access scope granted to a delegate.
type DelegationScope string

const (
	DelegationScopeRead  DelegationScope = "read"
	DelegationScopeWrite DelegationScope = "write"
	DelegationScopeAdmin DelegationScope = "admin"
)

// IsValidDelegationScope returns true if the scope is a known value.
func IsValidDelegationScope(s DelegationScope) bool {
	switch s {
	case DelegationScopeRead, DelegationScopeWrite, DelegationScopeAdmin:
		return true
	}
	return false
}

// SecretDelegation represents a delegated access grant for a secret.
type SecretDelegation struct {
	Mount       string          `json:"mount"`
	Path        string          `json:"path"`
	DelegatedTo string          `json:"delegated_to"`
	GrantedBy   string          `json:"granted_by"`
	Scope       DelegationScope `json:"scope"`
	ExpiresAt   time.Time       `json:"expires_at"`
	GrantedAt   time.Time       `json:"granted_at"`
	Note        string          `json:"note,omitempty"`
}

// FullPath returns the canonical path for the delegation target.
func (d SecretDelegation) FullPath() string {
	return fmt.Sprintf("%s/%s", d.Mount, d.Path)
}

// IsExpired returns true if the delegation has passed its expiry time.
func (d SecretDelegation) IsExpired() bool {
	if d.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(d.ExpiresAt)
}

// Validate checks that the delegation record is well-formed.
func (d SecretDelegation) Validate() error {
	if d.Mount == "" {
		return errors.New("delegation: mount is required")
	}
	if d.Path == "" {
		return errors.New("delegation: path is required")
	}
	if d.DelegatedTo == "" {
		return errors.New("delegation: delegated_to is required")
	}
	if d.GrantedBy == "" {
		return errors.New("delegation: granted_by is required")
	}
	if !IsValidDelegationScope(d.Scope) {
		return fmt.Errorf("delegation: unknown scope %q", d.Scope)
	}
	if !d.ExpiresAt.IsZero() && !d.ExpiresAt.After(time.Now()) {
		return errors.New("delegation: expires_at must be in the future")
	}
	return nil
}
