package vault

import (
	"testing"
	"time"
)

var baseCoverage = SecretCoverage{
	Mount:      "secret",
	Path:       "app/db",
	Status:     CoverageStatusFull,
	Score:      95,
	AssessedBy: "alice",
	AssessedAt: time.Now().UTC(),
}

func TestIsValidCoverageStatus_Known(t *testing.T) {
	for _, s := range []CoverageStatus{CoverageStatusFull, CoverageStatusPartial, CoverageStatusNone} {
		if !IsValidCoverageStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidCoverageStatus_Unknown(t *testing.T) {
	if IsValidCoverageStatus("unknown-level") {
		t.Error("expected unknown status to be invalid")
	}
}

func TestSecretCoverage_FullPath(t *testing.T) {
	if got := baseCoverage.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected FullPath: %s", got)
	}
}

func TestSecretCoverage_Validate_Valid(t *testing.T) {
	if err := baseCoverage.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretCoverage_Validate_MissingMount(t *testing.T) {
	c := baseCoverage
	c.Mount = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretCoverage_Validate_InvalidScore(t *testing.T) {
	c := baseCoverage
	c.Score = 150
	if err := c.Validate(); err == nil {
		t.Error("expected error for score > 100")
	}
}

func TestSecretCoverage_Validate_InvalidStatus(t *testing.T) {
	c := baseCoverage
	c.Status = "bogus"
	if err := c.Validate(); err == nil {
		t.Error("expected error for invalid status")
	}
}

func TestNewSecretCoverageRegistry_NotNil(t *testing.T) {
	if NewSecretCoverageRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestCoverageRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretCoverageRegistry()
	if err := r.Set(baseCoverage); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	got, ok := r.Get(baseCoverage.Mount, baseCoverage.Path)
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.Score != baseCoverage.Score {
		t.Errorf("score mismatch: got %d", got.Score)
	}
}

func TestCoverageRegistry_Set_SetsAssessedAt(t *testing.T) {
	r := NewSecretCoverageRegistry()
	c := baseCoverage
	c.AssessedAt = time.Time{}
	if err := r.Set(c); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	got, _ := r.Get(c.Mount, c.Path)
	if got.AssessedAt.IsZero() {
		t.Error("expected AssessedAt to be stamped")
	}
}

func TestCoverageRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretCoverageRegistry()
	_, ok := r.Get("missing", "path")
	if ok {
		t.Error("expected not found")
	}
}

func TestCoverageRegistry_Remove(t *testing.T) {
	r := NewSecretCoverageRegistry()
	_ = r.Set(baseCoverage)
	r.Remove(baseCoverage.Mount, baseCoverage.Path)
	_, ok := r.Get(baseCoverage.Mount, baseCoverage.Path)
	if ok {
		t.Error("expected record to be removed")
	}
}
