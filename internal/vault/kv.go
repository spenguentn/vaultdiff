package vault

// KVVersion represents the KV secrets engine version.
type KVVersion int

const (
	// KVv1 is the original KV secrets engine (no versioning).
	KVv1 KVVersion = 1
	// KVv2 is the versioned KV secrets engine.
	KVv2 KVVersion = 2
)

// KVConfig holds configuration for a KV secrets engine mount.
type KVConfig struct {
	// MountPath is the mount point of the KV engine (e.g. "secret").
	MountPath string
	// Version is the KV engine version (1 or 2).
	Version KVVersion
}

// Validate returns an error if the KVConfig is invalid.
func (c KVConfig) Validate() error {
	if c.MountPath == "" {
		return ErrMissingMountPath
	}
	if c.Version != KVv1 && c.Version != KVv2 {
		return ErrInvalidKVVersion
	}
	return nil
}

// IsVersioned reports whether the KV engine supports secret versioning.
func (c KVConfig) IsVersioned() bool {
	return c.Version == KVv2
}

// DataPrefix returns the path prefix used to read secret data.
// KVv2 engines use a "data/" infix; KVv1 engines do not.
func (c KVConfig) DataPrefix() string {
	if c.Version == KVv2 {
		return c.MountPath + "/data"
	}
	return c.MountPath
}

// MetadataPrefix returns the path prefix used to read secret metadata.
// Only meaningful for KVv2; returns empty string for KVv1.
func (c KVConfig) MetadataPrefix() string {
	if c.Version == KVv2 {
		return c.MountPath + "/metadata"
	}
	return ""
}

// KV-related sentinel errors.
var (
	ErrMissingMountPath = kvErr("mount path must not be empty")
	ErrInvalidKVVersion = kvErr("kv version must be 1 or 2")
)

type kvErr string

func (e kvErr) Error() string { return string(e) }
