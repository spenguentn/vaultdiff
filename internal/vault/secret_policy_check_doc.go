// Package vault provides policy check utilities for evaluating whether
// a given capability is permitted on a secret path.
//
// # Policy Check
//
// PolicyCheckRequest describes the target mount, path, and capability
// (e.g. "read", "write", "delete") to evaluate.
//
// CheckSecretPolicy matches the request against a slice of PolicyRule values
// and returns a PolicyCheckResult indicating whether the operation is allowed.
//
// Rule paths support trailing wildcard matching using "*", e.g.:
//
//	"secret/myapp/*" matches "secret/myapp/db" and "secret/myapp/cache"
//
// Example:
//
//	req := vault.PolicyCheckRequest{
//		Mount:      "secret",
//		Path:       "myapp/db",
//		Capability: "read",
//	}
//
//	result, err := vault.CheckSecretPolicy(req, rules)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if !result.IsAllowed() {
//		log.Printf("access denied: %s", result.Reason)
//	}
package vault
