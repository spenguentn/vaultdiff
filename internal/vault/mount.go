package vault

import (
	"errors"
	"strings"
)

// MountType represents the type of a Vault secrets mount.
type MountType string

const (
	MountTypeKV1    MountType = "kv-v1"
	MountTypeKV2    MountType = "kv-v2"
	MountTypePKI    MountType = "pki"
	MountTypeTransit MountType = "transit"
	MountTypeGeneric MountType = "generic"
)

// MountInfo describes a single secrets mount point in Vault.
type MountInfo struct {
	Path        string    `json:"path"`
	Type        MountType `json:"type"`
	Description string    `json:"description"`
	Versioned   bool      `json:"versioned"`
}

// Validate returns an error if the MountInfo is not usable.
func (m MountInfo) Validate() error {
	if strings.TrimSpace(m.Path) == "" {
		return errors.New("mount path must not be empty")
	}
	if strings.TrimSpace(string(m.Type)) == "" {
		return errors.New("mount type must not be empty")
	}
	return nil
}

// IsKV reports whether the mount is a KV secrets engine (v1 or v2).
func (m MountInfo) IsKV() bool {
	return m.Type == MountTypeKV1 || m.Type == MountTypeKV2
}

// ParseMountType normalises a raw string into a MountType.
// Unknown values are returned as-is cast to MountType.
func ParseMountType(raw string) MountType {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "kv", "kv-v1":
		return MountTypeKV1
	case "kv-v2":
		return MountTypeKV2
	case "pki":
		return MountTypePKI
	case "transit":
		return MountTypeTransit
	case "generic":
		return MountTypeGeneric
	default:
		return MountType(raw)
	}
}
