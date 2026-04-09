package vault

import "fmt"

// Environment represents a named Vault environment (e.g. staging, production)
// with its own address and mount path.
type Environment struct {
	Name      string
	Address   string
	MountPath string
	Token     string
}

// Validate checks that all required fields are set on the Environment.
func (e *Environment) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("environment name is required")
	}
	if e.Address == "" {
		return fmt.Errorf("environment %q: address is required", e.Name)
	}
	if e.MountPath == "" {
		return fmt.Errorf("environment %q: mount path is required", e.Name)
	}
	if e.Token == "" {
		return fmt.Errorf("environment %q: token is required", e.Name)
	}
	return nil
}

// SecretPath returns the full KV v2 path for a given secret key within
// this environment's mount.
func (e *Environment) SecretPath(key string) string {
	return fmt.Sprintf("%s/data/%s", e.MountPath, key)
}

// MetadataPath returns the full KV v2 metadata path for a given secret key.
func (e *Environment) MetadataPath(key string) string {
	return fmt.Sprintf("%s/metadata/%s", e.MountPath, key)
}
