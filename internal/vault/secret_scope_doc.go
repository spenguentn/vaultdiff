// Package vault provides secret scope management for vaultdiff.
//
// Secret scope defines the visibility boundary of a secret within an
// organisation. Supported scope levels are:
//
//   - local   – visible only within the owning service or component
//   - team    – shared across a single team's namespaces
//   - global  – accessible organisation-wide
//
// # Usage
//
//	registry := vault.NewSecretScopeRegistry()
//
//	scope := vault.SecretScope{
//		Mount: "secret",
//		Path:  "payments/api-key",
//		Level: vault.ScopeLevelTeam,
//		Owner: "team-payments",
//	}
//
//	if err := registry.Set(scope); err != nil {
//		log.Fatal(err)
//	}
//
//	entry, ok := registry.Get("secret", "payments/api-key")
//	if ok {
//		fmt.Println(entry.Level)
//	}
package vault
