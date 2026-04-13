// Package vault provides the SecretPin and SecretPinRegistry types for
// managing pinned secret versions within HashiCorp Vault.
//
// # Overview
//
// A SecretPin anchors a specific version of a secret so that automated
// rotation, promotion, or copy operations cannot overwrite it until the
// pin is explicitly removed or expires.
//
// # Usage
//
//	registry := vault.NewSecretPinRegistry()
//
//	pin := vault.SecretPin{
//		Mount:    "secret",
//		Path:     "app/database",
//		Version:  5,
//		PinnedBy: "deploy-bot",
//		Reason:   "v2.4.0 release freeze",
//	}
//
//	if err := registry.Pin(pin); err != nil {
//		log.Fatal(err)
//	}
//
//	if registry.IsPinned("secret", "app/database") {
//		fmt.Println("secret is pinned, skipping rotation")
//	}
//
//	registry.Unpin("secret", "app/database")
package vault
