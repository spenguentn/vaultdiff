package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// SecretChecksum holds a deterministic hash of a secret's key-value pairs.
type SecretChecksum struct {
	Mount     string    `json:"mount"`
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	Checksum  string    `json:"checksum"`
	ComputedAt time.Time `json:"computed_at"`
}

// FullPath returns the canonical mount+path string.
func (s SecretChecksum) FullPath() string {
	return fmt.Sprintf("%s/%s", strings.Trim(s.Mount, "/"), strings.Trim(s.Path, "/"))
}

// Matches reports whether two checksums are identical.
func (s SecretChecksum) Matches(other SecretChecksum) bool {
	return s.Checksum == other.Checksum
}

// Validate returns an error if the checksum is missing required fields.
func (s SecretChecksum) Validate() error {
	if s.Mount == "" {
		return fmt.Errorf("checksum: mount is required")
	}
	if s.Path == "" {
		return fmt.Errorf("checksum: path is required")
	}
	if s.Checksum == "" {
		return fmt.Errorf("checksum: checksum value is required")
	}
	return nil
}

// ComputeChecksum deterministically hashes the provided key-value data.
// Keys are sorted before hashing to ensure consistency regardless of map order.
func ComputeChecksum(mount, path string, version int, data map[string]interface{}) (SecretChecksum, error) {
	if mount == "" {
		return SecretChecksum{}, fmt.Errorf("compute checksum: mount is required")
	}
	if path == "" {
		return SecretChecksum{}, fmt.Errorf("compute checksum: path is required")
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%v;", k, data[k])
	}

	return SecretChecksum{
		Mount:      mount,
		Path:       path,
		Version:    version,
		Checksum:   hex.EncodeToString(h.Sum(nil)),
		ComputedAt: time.Now().UTC(),
	}, nil
}
