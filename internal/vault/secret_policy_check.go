package vault

import (
	"fmt"
	"strings"
)

// PolicyCheckRequest represents a request to verify a secret path against a policy.
type PolicyCheckRequest struct {
	Mount      string
	Path       string
	Capability string
}

// Validate returns an error if the request is incomplete.
func (r PolicyCheckRequest) Validate() error {
	if r.Mount == "" {
		return fmt.Errorf("mount is required")
	}
	if r.Path == "" {
		return fmt.Errorf("path is required")
	}
	if r.Capability == "" {
		return fmt.Errorf("capability is required")
	}
	return nil
}

// FullPath returns the combined mount and path.
func (r PolicyCheckRequest) FullPath() string {
	return strings.Trim(r.Mount, "/") + "/" + strings.Trim(r.Path, "/")
}

// PolicyCheckResult holds the outcome of a policy check.
type PolicyCheckResult struct {
	Request PolicyCheckRequest
	Allowed bool
	Reason  string
}

// IsAllowed returns true if the capability is permitted.
func (p PolicyCheckResult) IsAllowed() bool {
	return p.Allowed
}

// CheckSecretPolicy evaluates whether a given capability is allowed for a path
// based on the provided policy rules.
func CheckSecretPolicy(req PolicyCheckRequest, rules []PolicyRule) (PolicyCheckResult, error) {
	if err := req.Validate(); err != nil {
		return PolicyCheckResult{}, fmt.Errorf("invalid request: %w", err)
	}

	fullPath := req.FullPath()

	for _, rule := range rules {
		if ruleMatchesPath(rule.Path, fullPath) {
			if rule.HasCapability(req.Capability) {
				return PolicyCheckResult{
					Request: req,
					Allowed: true,
					Reason:  fmt.Sprintf("matched rule path %q", rule.Path),
				}, nil
			}
			return PolicyCheckResult{
				Request: req,
				Allowed: false,
				Reason:  fmt.Sprintf("rule %q does not grant %q", rule.Path, req.Capability),
			}, nil
		}
	}

	return PolicyCheckResult{
		Request: req,
		Allowed: false,
		Reason:  "no matching policy rule found",
	}, nil
}

// ruleMatchesPath checks if a policy rule path matches the given secret path.
// Supports trailing wildcard (*) matching.
func ruleMatchesPath(rulePath, secretPath string) bool {
	if strings.HasSuffix(rulePath, "*") {
		prefix := strings.TrimSuffix(rulePath, "*")
		return strings.HasPrefix(secretPath, prefix)
	}
	return rulePath == secretPath
}
