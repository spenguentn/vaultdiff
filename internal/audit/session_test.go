package audit_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/audit"
	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestNewSession_WithUser(t *testing.T) {
	s, err := audit.NewSession("staging", "secret/data/svc", 1, 2, "bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.User != "bob" {
		t.Errorf("expected user=bob, got %s", s.User)
	}
	if s.Environment != "staging" {
		t.Errorf("expected env=staging, got %s", s.Environment)
	}
}

func TestNewSession_AutoUser(t *testing.T) {
	s, err := audit.NewSession("dev", "secret/data/svc", 1, 2, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// User should be populated from the OS; just ensure it's non-empty.
	if s.User == "" {
		t.Error("expected non-empty user from OS, got empty string")
	}
}

func TestSession_BuildEntry_Summary(t *testing.T) {
	s, _ := audit.NewSession("prod", "secret/data/app", 5, 6, "carol")

	changes := []diff.Change{
		{Key: "A", Type: diff.Added},
		{Key: "B", Type: diff.Added},
		{Key: "C", Type: diff.Removed},
		{Key: "D", Type: diff.Modified},
		{Key: "E", Type: diff.Unchanged},
	}

	entry := s.BuildEntry(changes)

	if entry.Summary.Added != 2 {
		t.Errorf("expected added=2, got %d", entry.Summary.Added)
	}
	if entry.Summary.Removed != 1 {
		t.Errorf("expected removed=1, got %d", entry.Summary.Removed)
	}
	if entry.Summary.Modified != 1 {
		t.Errorf("expected modified=1, got %d", entry.Summary.Modified)
	}
	if entry.Summary.Unchanged != 1 {
		t.Errorf("expected unchanged=1, got %d", entry.Summary.Unchanged)
	}
	if entry.Path != "secret/data/app" {
		t.Errorf("expected path=secret/data/app, got %s", entry.Path)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
