package vault

import (
	"fmt"
	"strings"
)

// SecretPath represents a fully-qualified path to a KV v2 secret in Vault.
type SecretPath struct {
	Mount   string
	SubPath string
}

// NewSecretPath constructs a SecretPath from a mount and a sub-path,
// trimming any extraneous slashes.
func NewSecretPath(mount, subPath string) SecretPath {
	return SecretPath{
		Mount:   strings.Trim(mount, "/"),
		SubPath: strings.Trim(subPath, "/"),
	}
}

// DataPath returns the full API path used to read secret data.
// e.g. "secret/data/myapp/config"
func (p SecretPath) DataPath() string {
	return fmt.Sprintf("%s/data/%s", p.Mount, p.SubPath)
}

// MetadataPath returns the full API path used to read secret metadata.
// e.g. "secret/metadata/myapp/config"
func (p SecretPath) MetadataPath() string {
	return fmt.Sprintf("%s/metadata/%s", p.Mount, p.SubPath)
}

// String returns a human-readable representation of the path.
func (p SecretPath) String() string {
	return fmt.Sprintf("%s/%s", p.Mount, p.SubPath)
}

// Validate ensures neither Mount nor SubPath is empty.
func (p SecretPath) Validate() error {
	if p.Mount == "" {
		return fmt.Errorf("secret path: mount must not be empty")
	}
	if p.SubPath == "" {
		return fmt.Errorf("secret path: sub-path must not be empty")
	}
	return nil
}

// ParseSecretPath parses a raw string of the form "mount/sub/path" where the
// first segment is treated as the mount and the remainder as the sub-path.
func ParseSecretPath(raw string) (SecretPath, error) {
	raw = strings.Trim(raw, "/")
	idx := strings.Index(raw, "/")
	if idx < 0 {
		return SecretPath{}, fmt.Errorf("secret path %q must contain at least one '/' separating mount from sub-path", raw)
	}
	return NewSecretPath(raw[:idx], raw[idx+1:]), nil
}
