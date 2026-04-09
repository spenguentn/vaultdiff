package vault

import (
	"testing"
	"time"
)

func TestParseMetadata_FullData(t *testing.T) {
	data := map[string]interface{}{
		"current_version": json.Number("3"),
		"oldest_version":  json.Number("1"),
		"created_time":    "2024-01-01T00:00:00Z",
		"updated_time":    "2024-06-01T12:00:00Z",
		"versions": map[string]interface{}{
			"1": map[string]interface{}{
				"created_time":  "2024-01-01T00:00:00Z",
				"deletion_time": "",
				"destroyed":     false,
			},
			"2": map[string]interface{}{
				"created_time":  "2024-03-01T00:00:00Z",
				"deletion_time": "2024-04-01T00:00:00Z",
				"destroyed":     true,
			},
		},
	}

	meta, err := parseMetadata("secret/myapp", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.Path != "secret/myapp" {
		t.Errorf("expected path %q, got %q", "secret/myapp", meta.Path)
	}
	if meta.CurrentVersion != 3 {
		t.Errorf("expected current_version 3, got %d", meta.CurrentVersion)
	}
	if meta.OldestVersion != 1 {
		t.Errorf("expected oldest_version 1, got %d", meta.OldestVersion)
	}
	if len(meta.Versions) != 2 {
		t.Errorf("expected 2 versions, got %d", len(meta.Versions))
	}
	v2, ok := meta.Versions["2"]
	if !ok {
		t.Fatal("expected version 2 to exist")
	}
	if !v2.Destroyed {
		t.Error("expected version 2 to be destroyed")
	}
	if v2.DeletionTime.IsZero() {
		t.Error("expected deletion_time to be set for version 2")
	}
}

func TestParseMetadata_EmptyVersions(t *testing.T) {
	data := map[string]interface{}{
		"current_version": json.Number("0"),
		"oldest_version":  json.Number("0"),
		"created_time":    "2024-01-01T00:00:00Z",
		"updated_time":    "2024-01-01T00:00:00Z",
	}

	meta, err := parseMetadata("secret/empty", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(meta.Versions) != 0 {
		t.Errorf("expected 0 versions, got %d", len(meta.Versions))
	}
}

func TestParseMetadata_TimestampParsed(t *testing.T) {
	data := map[string]interface{}{
		"current_version": json.Number("1"),
		"oldest_version":  json.Number("1"),
		"created_time":    "2023-05-15T08:30:00Z",
		"updated_time":    "2023-05-15T08:30:00Z",
	}

	meta, err := parseMetadata("secret/ts", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2023, 5, 15, 8, 30, 0, 0, time.UTC)
	if !meta.CreatedTime.Equal(expected) {
		t.Errorf("expected created_time %v, got %v", expected, meta.CreatedTime)
	}
}
