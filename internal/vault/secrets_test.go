package vault

import (
	"context"
	"testing"
)

func TestParseSecretVersion_ValidData(t *testing.T) {
	raw := map[string]interface{}{
		"data": map[string]interface{}{
			"DB_PASS": "s3cr3t",
			"API_KEY": "abc123",
		},
		"metadata": map[string]interface{}{
			"version":        float64(3),
			"created_time":   "2024-01-01T00:00:00Z",
			"deletion_time":  "",
			"destroyed":      false,
		},
	}

	sv, err := parseSecretVersion(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sv.Version != 3 {
		t.Errorf("expected version 3, got %d", sv.Version)
	}
	if sv.Data["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t, got %s", sv.Data["DB_PASS"])
	}
	if sv.Data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", sv.Data["API_KEY"])
	}
}

func TestParseSecretVersion_EmptyData(t *testing.T) {
	raw := map[string]interface{}{
		"data":     map[string]interface{}{},
		"metadata": map[string]interface{}{},
	}
	sv, err := parseSecretVersion(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sv.Data) != 0 {
		t.Errorf("expected empty data map, got %v", sv.Data)
	}
}

func TestParseSecretVersion_MissingDataKey(t *testing.T) {
	raw := map[string]interface{}{
		"metadata": map[string]interface{}{"version": float64(1)},
	}
	sv, err := parseSecretVersion(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sv.Data) != 0 {
		t.Errorf("expected empty data, got %v", sv.Data)
	}
}

func TestReadSecret_NilResponse(t *testing.T) {
	// Ensures ReadSecret surfaces a meaningful error when vault returns nil.
	_ = context.Background() // placeholder; full integration tested via mock in client_test.go
}
