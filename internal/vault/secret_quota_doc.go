// Package vault provides secret quota management for vaultdiff.
//
// SecretQuota defines a rate limit on secret reads within a sliding time
// window. Quotas are scoped to a mount and optional path prefix.
//
// Usage:
//
//	registry := vault.NewSecretQuotaRegistry()
//	err := registry.Register(vault.SecretQuota{
//		Mount:      "secret",
//		Prefix:     "app/",
//		Scope:      vault.QuotaScopeMount,
//		MaxReads:   100,
//		WindowSize: time.Minute,
//	})
//
//	if registry.Allow("secret", "app/") {
//		// proceed with secret read
//	}
//
// Quotas are enforced using a sliding window counter. Accesses older than
// WindowSize are evicted on each call to Allow.
package vault
