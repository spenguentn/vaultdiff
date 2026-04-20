package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// IntegrityStatus represents the result of an integrity check.
type IntegrityStatus string

const (
	IntegrityStatusOK       IntegrityStatus = "ok"
	IntegrityStatusTampered IntegrityStatus = "tampered"
	IntegrityStatusUnknown  IntegrityStatus = "unknown"
)

// IsValidIntegrityStatus reports whether s is a known integrity status.
func IsValidIntegrityStatus(s IntegrityStatus) bool {
	switch s {
	case IntegrityStatusOK, IntegrityStatusTampered, IntegrityStatusUnknown:
		return true
	}
	return false
}

// SecretIntegrityRecord holds the integrity fingerprint for a secret.
type SecretIntegrityRecord struct {
	Mount       string          `json:"mount"`
	Path        string          `json:"path"`
	Fingerprint string          `json:"fingerprint"`
	Status      IntegrityStatus `json:"status"`
	CheckedAt   time.Time       `json:"checked_at"`
	CheckedBy   string          `json:"checked_by"`
}

// FullPath returns the canonical mount+path key.
func (r *SecretIntegrityRecord) FullPath() string {
	return fmt.Sprintf("%s/%s", strings.Trim(r.Mount, "/"), strings.Trim(r.Path, "/"))
}

// Validate returns an error if the record is missing required fields.
func (r *SecretIntegrityRecord) Validate() error {
	if r.Mount == "" {
		return errors.New("integrity: mount is required")
	}
	if r.Path == "" {
		return errors.New("integrity: path is required")
	}
	if r.Fingerprint == "" {
		return errors.New("integrity: fingerprint is required")
	}
	if !IsValidIntegrityStatus(r.Status) {
		return fmt.Errorf("integrity: unknown status %q", r.Status)
	}
	if r.CheckedBy == "" {
		return errors.New("integrity: checked_by is required")
	}
	return nil
}

// ComputeIntegrityFingerprint returns a deterministic SHA-256 hex digest of the
// provided secret data map.
func ComputeIntegrityFingerprint(data map[string]string) (string, error) {
	if data == nil {
		return "", errors.New("integrity: data must not be nil")
	}
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s;", k, data[k])
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
