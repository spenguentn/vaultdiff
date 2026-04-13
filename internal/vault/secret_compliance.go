package vault

import (
	"fmt"
	"time"
)

// ComplianceStatus represents the compliance state of a secret.
type ComplianceStatus string

const (
	ComplianceStatusCompliant    ComplianceStatus = "compliant"
	ComplianceStatusNonCompliant ComplianceStatus = "non_compliant"
	ComplianceStatusPending      ComplianceStatus = "pending"
	ComplianceStatusExempt       ComplianceStatus = "exempt"
)

// IsValidComplianceStatus returns true if the status is a known value.
func IsValidComplianceStatus(s ComplianceStatus) bool {
	switch s {
	case ComplianceStatusCompliant, ComplianceStatusNonCompliant,
		ComplianceStatusPending, ComplianceStatusExempt:
		return true
	}
	return false
}

// ComplianceRecord captures the compliance evaluation result for a secret.
type ComplianceRecord struct {
	Mount      string           `json:"mount"`
	Path       string           `json:"path"`
	Status     ComplianceStatus `json:"status"`
	Framework  string           `json:"framework"`
	Reason     string           `json:"reason,omitempty"`
	EvaluatedBy string          `json:"evaluated_by"`
	EvaluatedAt time.Time       `json:"evaluated_at"`
}

// FullPath returns the canonical mount+path identifier.
func (c *ComplianceRecord) FullPath() string {
	return fmt.Sprintf("%s/%s", c.Mount, c.Path)
}

// Validate returns an error if the record is missing required fields.
func (c *ComplianceRecord) Validate() error {
	if c.Mount == "" {
		return fmt.Errorf("compliance record: mount is required")
	}
	if c.Path == "" {
		return fmt.Errorf("compliance record: path is required")
	}
	if !IsValidComplianceStatus(c.Status) {
		return fmt.Errorf("compliance record: unknown status %q", c.Status)
	}
	if c.Framework == "" {
		return fmt.Errorf("compliance record: framework is required")
	}
	if c.EvaluatedBy == "" {
		return fmt.Errorf("compliance record: evaluated_by is required")
	}
	return nil
}

// complianceKey returns the registry key for a compliance record.
func complianceKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// ComplianceRegistry stores compliance records keyed by mount+path.
type ComplianceRegistry struct {
	records map[string]*ComplianceRecord
}

// NewComplianceRegistry returns an initialised ComplianceRegistry.
func NewComplianceRegistry() *ComplianceRegistry {
	return &ComplianceRegistry{records: make(map[string]*ComplianceRecord)}
}

// Record validates and stores a compliance record.
func (r *ComplianceRegistry) Record(rec *ComplianceRecord) error {
	if err := rec.Validate(); err != nil {
		return err
	}
	if rec.EvaluatedAt.IsZero() {
		rec.EvaluatedAt = time.Now().UTC()
	}
	r.records[complianceKey(rec.Mount, rec.Path)] = rec
	return nil
}

// Get returns the compliance record for the given mount and path.
func (r *ComplianceRegistry) Get(mount, path string) (*ComplianceRecord, bool) {
	v, ok := r.records[complianceKey(mount, path)]
	return v, ok
}

// Remove deletes the compliance record for the given mount and path.
func (r *ComplianceRegistry) Remove(mount, path string) {
	delete(r.records, complianceKey(mount, path))
}

// All returns every stored compliance record.
func (r *ComplianceRegistry) All() []*ComplianceRecord {
	out := make([]*ComplianceRecord, 0, len(r.records))
	for _, v := range r.records {
		out = append(out, v)
	}
	return out
}
