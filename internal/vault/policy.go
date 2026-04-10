package vault

import (
	"errors"
	"fmt"
	"strings"
)

// PolicyCapability represents a single Vault policy capability.
type PolicyCapability string

const (
	CapRead   PolicyCapability = "read"
	CapList   PolicyCapability = "list"
	CapCreate PolicyCapability = "create"
	CapUpdate PolicyCapability = "update"
	CapDelete PolicyCapability = "delete"
	CapDeny   PolicyCapability = "deny"
)

// PolicyRule represents a single path rule within a Vault policy.
type PolicyRule struct {
	Path         string
	Capabilities []PolicyCapability
}

// Validate returns an error if the rule is not well-formed.
func (r PolicyRule) Validate() error {
	if strings.TrimSpace(r.Path) == "" {
		return errors.New("policy rule path must not be empty")
	}
	if len(r.Capabilities) == 0 {
		return errors.New("policy rule must have at least one capability")
	}
	return nil
}

// HasCapability reports whether the rule grants the given capability.
func (r PolicyRule) HasCapability(c PolicyCapability) bool {
	for _, cap := range r.Capabilities {
		if cap == c {
			return true
		}
	}
	return false
}

// Policy represents a named Vault policy composed of path rules.
type Policy struct {
	Name  string
	Rules []PolicyRule
}

// Validate returns an error if the policy is not well-formed.
func (p Policy) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("policy name must not be empty")
	}
	for i, r := range p.Rules {
		if err := r.Validate(); err != nil {
			return fmt.Errorf("rule[%d]: %w", i, err)
		}
	}
	return nil
}

// RuleCount returns the number of rules in the policy.
func (p Policy) RuleCount() int {
	return len(p.Rules)
}

// AllowsRead reports whether any rule for the given path grants read access.
func (p Policy) AllowsRead(path string) bool {
	for _, r := range p.Rules {
		if r.Path == path && r.HasCapability(CapRead) {
			return true
		}
	}
	return false
}
