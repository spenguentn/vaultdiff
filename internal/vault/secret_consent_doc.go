// Package vault provides types and registries for managing HashiCorp Vault
// secrets, including access consent tracking.
//
// # Secret Consent
//
// SecretConsent represents an explicit, auditable approval for a named
// principal to access a specific secret path.  Each record captures who
// granted access, on whose behalf, the current status, and an optional
// expiry window.
//
// Supported statuses:
//
//	- pending  – consent has been requested but not yet decided
//	- granted  – access is approved
//	- revoked  – a previously granted consent has been withdrawn
//	- expired  – the consent window has lapsed
//
// Use SecretConsentRegistry to store, retrieve, revoke, and remove consent
// records in a concurrency-safe manner.
//
// Example:
//
//	reg := vault.NewSecretConsentRegistry()
//	err := reg.Grant(vault.SecretConsent{
//	    Mount:     "secret",
//	    Path:      "myapp/db",
//	    GrantedTo: "alice",
//	    GrantedBy: "admin",
//	    Status:    vault.ConsentGranted,
//	    GrantedAt: time.Now(),
//	})
package vault
