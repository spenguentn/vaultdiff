package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/output"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", ChangeType: diff.Added, OldValue: "", NewValue: "localhost"},
		{Key: "DB_PASS", ChangeType: diff.Modified, OldValue: "old", NewValue: "new"},
		{Key: "API_KEY", ChangeType: diff.Unchanged, OldValue: "abc", NewValue: "abc"},
	}
}

func TestWriter_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(output.FormatText, nil, &buf)
	if err := w.Write(sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in text output, got: %s", out)
	}
}

func TestWriter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(output.FormatJSON, nil, &buf)
	if err := w.Write(sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(parsed) != 3 {
		t.Errorf("expected 3 results, got %d", len(parsed))
	}
	if parsed[0]["key"] != "DB_HOST" {
		t.Errorf("expected first key to be DB_HOST, got %v", parsed[0]["key"])
	}
}

func TestWriter_JSONFormat_MasksSecrets(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(output.FormatJSON, []string{"DB_PASS"}, &buf)
	if err := w.Write(sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	for _, item := range parsed {
		if item["key"] == "DB_PASS" {
			if item["new_value"] == "new" {
				t.Errorf("expected DB_PASS new_value to be masked, got 'new'")
			}
			return
		}
	}
	t.Error("DB_PASS not found in output")
}

func TestWriter_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(output.Format("xml"), nil, &buf)
	if err := w.Write(sampleResults()); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestNewWriter_DefaultsToStdout(t *testing.T) {
	w := output.NewWriter(output.FormatText, nil, nil)
	if w == nil {
		t.Error("expected non-nil writer")
	}
}
