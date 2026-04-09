package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultdiff/internal/audit"
	"github.com/yourusername/vaultdiff/internal/diff"
)

func sampleEntry() audit.Entry {
	return audit.Entry{
		Timestamp:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Environment: "production",
		Path:        "secret/data/app/config",
		FromVersion: 3,
		ToVersion:   4,
		User:        "alice",
		Changes: []diff.Change{
			{Key: "DB_PASSWORD", Type: diff.Modified, OldValue: "old", NewValue: "new"},
		},
		Summary: diff.Summary{Added: 0, Removed: 0, Modified: 1},
	}
}

func TestLogger_WriteJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf, audit.FormatJSON)

	if err := logger.Write(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}

	if got.Environment != "production" {
		t.Errorf("expected environment=production, got %s", got.Environment)
	}
	if got.Summary.Modified != 1 {
		t.Errorf("expected modified=1, got %d", got.Summary.Modified)
	}
}

func TestLogger_WriteText(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf, audit.FormatText)

	if err := logger.Write(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := buf.String()
	if !strings.Contains(line, "env=production") {
		t.Errorf("expected env=production in output, got: %s", line)
	}
	if !strings.Contains(line, "versions=3->4") {
		t.Errorf("expected versions=3->4 in output, got: %s", line)
	}
	if !strings.Contains(line, "user=alice") {
		t.Errorf("expected user=alice in output, got: %s", line)
	}
}

func TestLogger_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf, audit.Format("xml"))

	if err := logger.Write(sampleEntry()); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
