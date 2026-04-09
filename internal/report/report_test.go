package report_test

import (
	"testing"
	"time"

	"github.com/vaultdiff/internal/audit"
	"github.com/vaultdiff/internal/diff"
	"github.com/vaultdiff/internal/report"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_PASS", ChangeType: diff.Added, NewValue: "secret"},
		{Key: "API_KEY", ChangeType: diff.Modified, OldValue: "old", NewValue: "new"},
		{Key: "HOST", ChangeType: diff.Unchanged, OldValue: "localhost", NewValue: "localhost"},
	}
}

func TestNew_FieldsSet(t *testing.T) {
	session := audit.NewSession("tester")
	results := sampleResults()
	before := time.Now().UTC()
	rep := report.New(session, results, "secret/dev", "secret/prod")
	after := time.Now().UTC()

	if rep.Session != session {
		t.Error("expected session to be set")
	}
	if len(rep.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(rep.Results))
	}
	if rep.SourcePath != "secret/dev" {
		t.Errorf("unexpected source path: %s", rep.SourcePath)
	}
	if rep.TargetPath != "secret/prod" {
		t.Errorf("unexpected target path: %s", rep.TargetPath)
	}
	if rep.GeneratedAt.Before(before) || rep.GeneratedAt.After(after) {
		t.Error("GeneratedAt timestamp out of expected range")
	}
}

func TestReport_Summary(t *testing.T) {
	session := audit.NewSession("")
	rep := report.New(session, sampleResults(), "a", "b")
	s := rep.Summary()
	if s.Added != 1 {
		t.Errorf("expected 1 added, got %d", s.Added)
	}
	if s.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", s.Modified)
	}
	if s.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", s.Unchanged)
	}
}
