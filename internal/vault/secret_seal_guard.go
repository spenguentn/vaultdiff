package vault

import (
	"errors"
	"fmt"
	"time"
)

// SealGuardAction defines what action triggered the seal guard.
type SealGuardAction string

const (
	SealGuardActionRead   SealGuardAction = "read"
	SealGuardActionWrite  SealGuardAction = "write"
	SealGuardActionDelete SealGuardAction = "delete"
)

// SealGuardEntry records a blocked operation due to a sealed or unhealthy vault.
type SealGuardEntry struct {
	Mount     string
	Path      string
	Action    SealGuardAction
	BlockedAt time.Time
	Reason    string
}

// FullPath returns the mount-prefixed path for the guarded secret.
func (e SealGuardEntry) FullPath() string {
	return fmt.Sprintf("%s/%s", e.Mount, e.Path)
}

// Validate checks that the SealGuardEntry has required fields.
func (e SealGuardEntry) Validate() error {
	if e.Mount == "" {
		return errors.New("seal guard: mount is required")
	}
	if e.Path == "" {
		return errors.New("seal guard: path is required")
	}
	if e.Action == "" {
		return errors.New("seal guard: action is required")
	}
	if e.BlockedAt.IsZero() {
		return errors.New("seal guard: blocked_at is required")
	}
	return nil
}

// IsValidSealGuardAction reports whether the given action string is known.
func IsValidSealGuardAction(a SealGuardAction) bool {
	switch a {
	case SealGuardActionRead, SealGuardActionWrite, SealGuardActionDelete:
		return true
	}
	return false
}

// SealGuardRegistry stores blocked operation entries keyed by mount+path.
type SealGuardRegistry struct {
	entries map[string][]SealGuardEntry
}

func guardKey(mount, path string) string {
	return mount + "/" + path
}

// NewSealGuardRegistry returns an initialised SealGuardRegistry.
func NewSealGuardRegistry() *SealGuardRegistry {
	return &SealGuardRegistry{entries: make(map[string][]SealGuardEntry)}
}

// Record appends a validated SealGuardEntry to the registry.
func (r *SealGuardRegistry) Record(e SealGuardEntry) error {
	if e.BlockedAt.IsZero() {
		e.BlockedAt = time.Now().UTC()
	}
	if err := e.Validate(); err != nil {
		return err
	}
	k := guardKey(e.Mount, e.Path)
	r.entries[k] = append(r.entries[k], e)
	return nil
}

// Get returns all guard entries for the given mount and path.
func (r *SealGuardRegistry) Get(mount, path string) ([]SealGuardEntry, bool) {
	v, ok := r.entries[guardKey(mount, path)]
	return v, ok
}

// Count returns the total number of recorded guard entries across all paths.
func (r *SealGuardRegistry) Count() int {
	n := 0
	for _, v := range r.entries {
		n += len(v)
	}
	return n
}
