// Package vault provides expiry tracking for Vault secrets.
//
// # Secret Expiry
//
// ExpiryPolicy records when a secret should be considered expired and how far
// in advance operators should be warned. Policies are stored in a
// SecretExpiryRegistry and can be evaluated in bulk via CheckAll.
//
// Usage:
//
//	reg := vault.NewSecretExpiryRegistry()
//
//	policy := vault.ExpiryPolicy{
//		Mount:      "secret",
//		Path:       "prod/db-password",
//		ExpiresAt:  time.Now().Add(30 * 24 * time.Hour),
//		WarnBefore: 7 * 24 * time.Hour,
//		Owner:      "team-infra",
//	}
//
//	if err := reg.Register(policy); err != nil {
//		log.Fatal(err)
//	}
//
//	for _, status := range reg.CheckAll(time.Now()) {
//		fmt.Printf("%s: %s\n", status.Policy.FullPath(), status)
//	}
package vault
