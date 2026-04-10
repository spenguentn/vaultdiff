package vault

import (
	"bytes"
	"strings"
	"testing"
)

func sampleRecord() ExportRecord {
	e := NewExporter(ExportOptions{})
	return e.Export("secret/svc", 1, map[string]string{
		"TOKEN": "abc123",
	})
}

func TestExportWriter_JSON(t *testing.T) {
	var buf bytes.Buffer
	ew := NewExportWriter(&buf, ExportOptions{Format: ExportFormatJSON})
	if err := ew.Write(sampleRecord()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"path"`) {
		t.Error("JSON output missing path field")
	}
}

func TestExportWriter_Env(t *testing.T) {
	var buf bytes.Buffer
	ew := NewExportWriter(&buf, ExportOptions{Format: ExportFormatEnv})
	if err := ew.Write(sampleRecord()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "TOKEN=") {
		t.Error("env output missing TOKEN key")
	}
}

func TestExportWriter_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	ew := NewExportWriter(&buf, ExportOptions{Format: "yaml"})
	err := ew.Write(sampleRecord())
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExportWriter_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	ew := NewExportWriter(&buf, ExportOptions{})
	if ew.opts.Format != ExportFormatJSON {
		t.Errorf("expected default JSON, got %q", ew.opts.Format)
	}
}

func TestExportWriter_WriteAll(t *testing.T) {
	var buf bytes.Buffer
	ew := NewExportWriter(&buf, ExportOptions{Format: ExportFormatEnv})
	records := []ExportRecord{sampleRecord(), sampleRecord()}
	if err := ew.WriteAll(records); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	count := strings.Count(buf.String(), "TOKEN=")
	if count != 2 {
		t.Errorf("expected 2 TOKEN lines, got %d", count)
	}
}
