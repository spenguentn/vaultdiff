package vault

import (
	"fmt"
	"time"
)

// TenureStatus represents how long a secret has been active.
type TenureStatus string

const (
	TenureNew     TenureStatus = "new"
	TenureActive  TenureStatus = "active"
	TenureMature  TenureStatus = "mature"
	TenureVeteran TenureStatus = "veteran"
)

// IsValidTenureStatus returns true if the given status is a known tenure status.
func IsValidTenureStatus(s TenureStatus) bool {
	switch s {
	case TenureNew, TenureActive, TenureMature, TenureVeteran:
		return true
	}
	return false
}

// SecretTenure records how long a secret has existed and its derived status.
type SecretTenure struct {
	Mount     string       `json:"mount"`
	Path      string       `json:"path"`
	CreatedAt time.Time    `json:"created_at"`
	Status    TenureStatus `json:"status"`
}

// FullPath returns the canonical mount+path string.
func (t SecretTenure) FullPath() string {
	return fmt.Sprintf("%s/%s", t.Mount, t.Path)
}

// AgeDays returns the number of full days since the secret was created.
func (t SecretTenure) AgeDays() int {
	return int(time.Since(t.CreatedAt).Hours() / 24)
}

// Validate returns an error if the tenure record is incomplete.
func (t SecretTenure) Validate() error {
	if t.Mount == "" {
		return fmt.Errorf("tenure: mount is required")
	}
	if t.Path == "" {
		return fmt.Errorf("tenure: path is required")
	}
	if t.CreatedAt.IsZero() {
		return fmt.Errorf("tenure: created_at is required")
	}
	if !IsValidTenureStatus(t.Status) {
		return fmt.Errorf("tenure: unknown status %q", t.Status)
	}
	return nil
}

// ComputeTenure derives a SecretTenure from a creation timestamp.
func ComputeTenure(mount, path string, createdAt time.Time) SecretTenure {
	days := int(time.Since(createdAt).Hours() / 24)
	var status TenureStatus
	switch {
	case days < 30:
		status = TenureNew
	case days < 180:
		status = TenureActive
	case days < 365:
		status = TenureMature
	default:
		status = TenureVeteran
	}
	return SecretTenure{
		Mount:     mount,
		Path:      path,
		CreatedAt: createdAt,
		Status:    status,
	}
}
