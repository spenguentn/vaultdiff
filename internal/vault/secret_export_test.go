package vault

import (
	"encoding/json"
	"strings"
	"testing"
)

func baseExportData() map[string]string {
	return map[string]string{
		"DB_PASS": "secret123",
		"API_KEY": "key-abc",
	}
}

func TestNewExporter_DefaultFormat(t *testing.T) {
	e := NewExporter(ExportOptions{})
	if e.opts.Format != ExportFormatJSON {
		t.Fatalf("expected JSON format, got %q", e.opts.Format)
	}
}

func TestNewExporter_CustomFormat(t *testing.T) {
	e := NewExporter(ExportOptions{Format: ExportFormatEnv})
	if e.opts.Format != ExportFormatEnv {
		t.Fatalf("expected env format, got %q", e.opts.Format)
	}
}

func TestExporter_Export_PreservesData(t *testing.T) {
	e := NewExporter(ExportOptions{})
	rec := e.Export("secret/app", 3, baseExportData())
	if rec.Path != "secret/app" {
		t.Errorf("path mismatch: %q", rec.Path)
	}
	if rec.Version != 3 {
		t.Errorf("version mismatch: %d", rec.Version)
	}
	if rec.Data["DB_PASS"] != "secret123" {
		t.Errorf("unexpected value: %q", rec.Data["DB_PASS"])
	}
}

func TestExporter_Export_MasksValues(t *testing.T) {
	e := NewExporter(ExportOptions{MaskValues: true})
	rec := e.Export("secret/app", 1, baseExportData())
	for k, v := range rec.Data {
		if v != "***" {
			t.Errorf("key %q not masked: %q", k, v)
		}
	}
}

func TestMarshalExportRecord_ValidJSON(t *testing.T) {
	e := NewExporter(ExportOptions{})
	rec := e.Export("secret/app", 2, baseExportData())
	b, err := MarshalExportRecord(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["path"] != "secret/app" {
		t.Errorf("path not in JSON output")
	}
}

func TestMarshalEnv_ContainsKeys(t *testing.T) {
	e := NewExporter(ExportOptions{})
	rec := e.Export("secret/app", 1, baseExportData())
	raw := string(MarshalEnv(rec))
	if !strings.Contains(raw, "DB_PASS=") {
		t.Errorf("env output missing DB_PASS")
	}
	if !strings.Contains(raw, "API_KEY=") {
		t.Errorf("env output missing API_KEY")
	}
}

func TestExporter_Export_ExportedAtSet(t *testing.T) {
	e := NewExporter(ExportOptions{})
	rec := e.Export("secret/db", 1, baseExportData())
	if rec.ExportedAt.IsZero() {
		t.Error("ExportedAt should not be zero")
	}
}
