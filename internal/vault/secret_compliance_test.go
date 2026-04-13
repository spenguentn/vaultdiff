package vault

import (
	"testing"
	"time"
)

var baseCompliance = &ComplianceRecord{
	Mount:       "secret",
	Path:        "app/db",
	Status:      ComplianceStatusCompliant,
	Framework:   "SOC2",
	EvaluatedBy: "ci-bot",
}

func TestIsValidComplianceStatus_Known(t *testing.T) {
	for _, s := range []ComplianceStatus{
		ComplianceStatusCompliant, ComplianceStatusNonCompliant,
		ComplianceStatusPending, ComplianceStatusExempt,
	} {
		if !IsValidComplianceStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidComplianceStatus_Unknown(t *testing.T) {
	if IsValidComplianceStatus(ComplianceStatus("bogus")) {
		t.Error("expected bogus to be invalid")
	}
}

func TestComplianceRecord_FullPath(t *testing.T) {
	if got := baseCompliance.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestComplianceRecord_Validate_Valid(t *testing.T) {
	if err := baseCompliance.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestComplianceRecord_Validate_MissingMount(t *testing.T) {
	r := *baseCompliance
	r.Mount = ""
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestComplianceRecord_Validate_MissingPath(t *testing.T) {
	r := *baseCompliance
	r.Path = ""
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestComplianceRecord_Validate_UnknownStatus(t *testing.T) {
	r := *baseCompliance
	r.Status = "invalid"
	if err := r.Validate(); err == nil {
		t.Error("expected error for unknown status")
	}
}

func TestComplianceRecord_Validate_MissingFramework(t *testing.T) {
	r := *baseCompliance
	r.Framework = ""
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing framework")
	}
}

func TestNewComplianceRegistry_NotNil(t *testing.T) {
	if NewComplianceRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestComplianceRegistry_Record_And_Get(t *testing.T) {
	reg := NewComplianceRegistry()
	rec := *baseCompliance
	if err := reg.Record(&rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("secret", "app/db")
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.Framework != "SOC2" {
		t.Errorf("unexpected framework: %s", got.Framework)
	}
}

func TestComplianceRegistry_Record_SetsEvaluatedAt(t *testing.T) {
	reg := NewComplianceRegistry()
	rec := *baseCompliance
	rec.EvaluatedAt = time.Time{}
	_ = reg.Record(&rec)
	if rec.EvaluatedAt.IsZero() {
		t.Error("expected EvaluatedAt to be set")
	}
}

func TestComplianceRegistry_Record_Invalid(t *testing.T) {
	reg := NewComplianceRegistry()
	rec := *baseCompliance
	rec.Mount = ""
	if err := reg.Record(&rec); err == nil {
		t.Error("expected error for invalid record")
	}
}

func TestComplianceRegistry_Get_NotFound(t *testing.T) {
	reg := NewComplianceRegistry()
	_, ok := reg.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestComplianceRegistry_Remove(t *testing.T) {
	reg := NewComplianceRegistry()
	rec := *baseCompliance
	_ = reg.Record(&rec)
	reg.Remove("secret", "app/db")
	_, ok := reg.Get("secret", "app/db")
	if ok {
		t.Error("expected record to be removed")
	}
}

func TestComplianceRegistry_All(t *testing.T) {
	reg := NewComplianceRegistry()
	r1 := *baseCompliance
	r2 := *baseCompliance
	r2.Path = "app/cache"
	_ = reg.Record(&r1)
	_ = reg.Record(&r2)
	if len(reg.All()) != 2 {
		t.Errorf("expected 2 records, got %d", len(reg.All()))
	}
}
