package vault

import (
	"errors"
	"fmt"
	"time"
)

// ProvenanceSource identifies where a secret originated.
type ProvenanceSource string

const (
	ProvenanceSourceManual    ProvenanceSource = "manual"
	ProvenanceSourceGenerated ProvenanceSource = "generated"
	ProvenanceSourceImported  ProvenanceSource = "imported"
	ProvenanceSourceReplicated ProvenanceSource = "replicated"
	ProvenanceSourceMigrated  ProvenanceSource = "migrated"
)

// IsValidProvenanceSource returns true if s is a known provenance source.
func IsValidProvenanceSource(s ProvenanceSource) bool {
	switch s {
	case ProvenanceSourceManual,
		ProvenanceSourceGenerated,
		ProvenanceSourceImported,
		ProvenanceSourceReplicated,
		ProvenanceSourceMigrated:
		return true
	}
	return false
}

// SecretProvenance records the origin and chain of custody for a secret.
type SecretProvenance struct {
	// Mount is the KV mount path.
	Mount string `json:"mount"`
	// Path is the secret path within the mount.
	Path string `json:"path"`
	// Source describes how the secret came to exist.
	Source ProvenanceSource `json:"source"`
	// CreatedBy is the actor who created or introduced the secret.
	CreatedBy string `json:"created_by"`
	// OriginMount is the original mount if the secret was copied or migrated.
	OriginMount string `json:"origin_mount,omitempty"`
	// OriginPath is the original path if the secret was copied or migrated.
	OriginPath string `json:"origin_path,omitempty"`
	// RecordedAt is when the provenance entry was created.
	RecordedAt time.Time `json:"recorded_at"`
	// Notes is an optional free-text explanation.
	Notes string `json:"notes,omitempty"`
}

// FullPath returns the canonical mount+path identifier.
func (p *SecretProvenance) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Mount, p.Path)
}

// Validate checks that the provenance record is well-formed.
func (p *SecretProvenance) Validate() error {
	if p.Mount == "" {
		return errors.New("provenance: mount is required")
	}
	if p.Path == "" {
		return errors.New("provenance: path is required")
	}
	if p.CreatedBy == "" {
		return errors.New("provenance: created_by is required")
	}
	if !IsValidProvenanceSource(p.Source) {
		return fmt.Errorf("provenance: unknown source %q", p.Source)
	}
	return nil
}

// provenanceKey builds the registry lookup key.
func provenanceKey(mount, path string) string {
	return mount + "|" + path
}

// SecretProvenanceRegistry stores provenance records keyed by mount+path.
type SecretProvenanceRegistry struct {
	records map[string]*SecretProvenance
}

// NewSecretProvenanceRegistry returns an initialised registry.
func NewSecretProvenanceRegistry() *SecretProvenanceRegistry {
	return &SecretProvenanceRegistry{
		records: make(map[string]*SecretProvenance),
	}
}

// Record stores a provenance entry, stamping RecordedAt if zero.
func (r *SecretProvenanceRegistry) Record(p *SecretProvenance) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if p.RecordedAt.IsZero() {
		p.RecordedAt = time.Now().UTC()
	}
	r.records[provenanceKey(p.Mount, p.Path)] = p
	return nil
}

// Get retrieves the provenance record for the given mount and path.
func (r *SecretProvenanceRegistry) Get(mount, path string) (*SecretProvenance, bool) {
	p, ok := r.records[provenanceKey(mount, path)]
	return p, ok
}

// Remove deletes a provenance record from the registry.
func (r *SecretProvenanceRegistry) Remove(mount, path string) {
	delete(r.records, provenanceKey(mount, path))
}

// All returns every provenance record currently stored.
func (r *SecretProvenanceRegistry) All() []*SecretProvenance {
	out := make([]*SecretProvenance, 0, len(r.records))
	for _, v := range r.records {
		out = append(out, v)
	}
	return out
}
