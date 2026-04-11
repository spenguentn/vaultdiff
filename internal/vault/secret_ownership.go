package vault

import (
	"errors"
	"fmt"
	"time"
)

// OwnershipRecord represents the ownership metadata for a secret.
type OwnershipRecord struct {
	Mount     string    `json:"mount"`
	Path      string    `json:"path"`
	Owner     string    `json:"owner"`
	Team      string    `json:"team,omitempty"`
	Contact   string    `json:"contact,omitempty"`
	AssignedAt time.Time `json:"assigned_at"`
}

// FullPath returns the canonical mount+path identifier.
func (o OwnershipRecord) FullPath() string {
	return fmt.Sprintf("%s/%s", o.Mount, o.Path)
}

// Validate checks that required fields are present.
func (o OwnershipRecord) Validate() error {
	if o.Mount == "" {
		return errors.New("ownership: mount is required")
	}
	if o.Path == "" {
		return errors.New("ownership: path is required")
	}
	if o.Owner == "" {
		return errors.New("ownership: owner is required")
	}
	return nil
}

// OwnershipRegistry stores and retrieves ownership records keyed by mount+path.
type OwnershipRegistry struct {
	records map[string]OwnershipRecord
}

// NewOwnershipRegistry creates an empty OwnershipRegistry.
func NewOwnershipRegistry() *OwnershipRegistry {
	return &OwnershipRegistry{
		records: make(map[string]OwnershipRecord),
	}
}

func ownershipKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// Assign stores an ownership record, overwriting any existing entry.
func (r *OwnershipRegistry) Assign(rec OwnershipRecord) error {
	if err := rec.Validate(); err != nil {
		return err
	}
	if rec.AssignedAt.IsZero() {
		rec.AssignedAt = time.Now().UTC()
	}
	r.records[ownershipKey(rec.Mount, rec.Path)] = rec
	return nil
}

// Get retrieves the ownership record for the given mount and path.
func (r *OwnershipRegistry) Get(mount, path string) (OwnershipRecord, bool) {
	rec, ok := r.records[ownershipKey(mount, path)]
	return rec, ok
}

// Remove deletes an ownership record. Returns false if not found.
func (r *OwnershipRegistry) Remove(mount, path string) bool {
	key := ownershipKey(mount, path)
	if _, ok := r.records[key]; !ok {
		return false
	}
	delete(r.records, key)
	return true
}

// All returns a slice of all registered ownership records.
func (r *OwnershipRegistry) All() []OwnershipRecord {
	out := make([]OwnershipRecord, 0, len(r.records))
	for _, rec := range r.records {
		out = append(out, rec)
	}
	return out
}
