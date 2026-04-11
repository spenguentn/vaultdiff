package vault

import (
	"fmt"
	"time"
)

// SealStatus represents the current seal state of a Vault instance.
type SealStatus int

const (
	SealStatusUnknown SealStatus = iota
	SealStatusSealed
	SealStatusUnsealed
)

// SealInfo holds the seal/unseal state returned by the Vault API.
type SealInfo struct {
	Sealed      bool      `json:"sealed"`
	Initialized bool      `json:"initialized"`
	Progress    int       `json:"progress"`
	Threshold   int       `json:"t"`
	Shares      int       `json:"n"`
	Version     string    `json:"version"`
	CheckedAt   time.Time `json:"-"`
}

// Status returns the SealStatus enum value.
func (s SealInfo) Status() SealStatus {
	if !s.Initialized {
		return SealStatusUnknown
	}
	if s.Sealed {
		return SealStatusSealed
	}
	return SealStatusUnsealed
}

// String returns a human-readable description of the seal state.
func (s SealInfo) String() string {
	switch s.Status() {
	case SealStatusSealed:
		return fmt.Sprintf("sealed (progress %d/%d)", s.Progress, s.Threshold)
	case SealStatusUnsealed:
		return fmt.Sprintf("unsealed (version %s)", s.Version)
	default:
		return "unknown (not initialized)"
	}
}

// IsSealed returns true when the Vault instance is sealed.
func (s SealInfo) IsSealed() bool {
	return s.Sealed
}

// UnsealProgress returns the fraction of unseal keys provided so far as a
// value between 0.0 and 1.0. Returns 0 if the threshold is not set or the
// instance is already unsealed.
func (s SealInfo) UnsealProgress() float64 {
	if s.Threshold <= 0 || !s.Sealed {
		return 0
	}
	return float64(s.Progress) / float64(s.Threshold)
}

// ParseSealInfo constructs a SealInfo from a raw Vault API response map.
func ParseSealInfo(raw map[string]any) (SealInfo, error) {
	if raw == nil {
		return SealInfo{}, fmt.Errorf("seal: nil response")
	}
	info := SealInfo{
		CheckedAt: time.Now().UTC(),
	}
	if v, ok := raw["sealed"].(bool); ok {
		info.Sealed = v
	}
	if v, ok := raw["initialized"].(bool); ok {
		info.Initialized = v
	}
	if v, ok := raw["progress"].(float64); ok {
		info.Progress = int(v)
	}
	if v, ok := raw["t"].(float64); ok {
		info.Threshold = int(v)
	}
	if v, ok := raw["n"].(float64); ok {
		info.Shares = int(v)
	}
	if v, ok := raw["version"].(string); ok {
		info.Version = v
	}
	return info, nil
}
