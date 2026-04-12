package vault

import (
	"testing"
)

func TestIsValidClassification_Known(t *testing.T) {
	for _, c := range []SecretClassification{
		ClassificationPublic,
		ClassificationInternal,
		ClassificationConfidential,
		ClassificationRestricted,
	} {
		if !IsValidClassification(c) {
			t.Errorf("expected %q to be valid", c)
		}
	}
}

func TestIsValidClassification_Unknown(t *testing.T) {
	if IsValidClassification("top-secret") {
		t.Error("expected unknown classification to be invalid")
	}
}

func TestClassifyRequest_Validate_Valid(t *testing.T) {
	req := ClassifyRequest{
		Mount:          "secret",
		Path:           "app/db",
		Classification: ClassificationConfidential,
		ClassifiedBy:   "alice",
	}
	if err := req.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClassifyRequest_Validate_MissingMount(t *testing.T) {
	req := ClassifyRequest{Path: "app/db", Classification: ClassificationPublic, ClassifiedBy: "alice"}
	if err := req.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestClassifyRequest_Validate_MissingPath(t *testing.T) {
	req := ClassifyRequest{Mount: "secret", Classification: ClassificationPublic, ClassifiedBy: "alice"}
	if err := req.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestClassifyRequest_Validate_InvalidClassification(t *testing.T) {
	req := ClassifyRequest{Mount: "secret", Path: "app/db", Classification: "unknown", ClassifiedBy: "alice"}
	if err := req.Validate(); err == nil {
		t.Error("expected error for invalid classification")
	}
}

func TestClassifyRequest_Validate_MissingClassifiedBy(t *testing.T) {
	req := ClassifyRequest{Mount: "secret", Path: "app/db", Classification: ClassificationRestricted}
	if err := req.Validate(); err == nil {
		t.Error("expected error for missing classifiedBy")
	}
}

func TestClassifyRequest_FullPath(t *testing.T) {
	req := ClassifyRequest{Mount: "secret/", Path: "/app/db"}
	if got := req.FullPath(); got != "secret/app/db" {
		t.Errorf("expected secret/app/db, got %q", got)
	}
}

func TestNewClassificationRecord_Valid(t *testing.T) {
	req := ClassifyRequest{
		Mount:          "secret",
		Path:           "app/key",
		Classification: ClassificationInternal,
		ClassifiedBy:   "bob",
	}
	rec, err := NewClassificationRecord(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Label != "internal" {
		t.Errorf("expected label 'internal', got %q", rec.Label)
	}
}

func TestNewClassificationRecord_Invalid(t *testing.T) {
	req := ClassifyRequest{}
	if _, err := NewClassificationRecord(req); err == nil {
		t.Error("expected error for invalid request")
	}
}
