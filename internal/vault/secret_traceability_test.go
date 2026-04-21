package vault

import (
	"testing"
	"time"
)

func TestIsValidTraceabilitySource_Known(t *testing.T) {
	for _, s := range []TraceabilitySource{
		TraceSourceManual, TraceSourcePipeline, TraceSourceImport,
		TraceSourceGenerated, TraceSourceMigrated,
	} {
		if !IsValidTraceabilitySource(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidTraceabilitySource_Unknown(t *testing.T) {
	if IsValidTraceabilitySource("unknown") {
		t.Error("expected 'unknown' to be invalid")
	}
}

func TestSecretTraceability_FullPath(t *testing.T) {
	tr := &SecretTraceability{Mount: "secret", Path: "app/db"}
	if got := tr.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretTraceability_Validate_Valid(t *testing.T) {
	tr := &SecretTraceability{
		Mount:    "secret",
		Path:     "app/db",
		Version:  1,
		Source:   TraceSourcePipeline,
		TracedBy: "ci-bot",
	}
	if err := tr.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretTraceability_Validate_MissingMount(t *testing.T) {
	tr := &SecretTraceability{Path: "app/db", Version: 1, Source: TraceSourceManual, TracedBy: "alice"}
	if err := tr.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretTraceability_Validate_ZeroVersion(t *testing.T) {
	tr := &SecretTraceability{Mount: "secret", Path: "app/db", Version: 0, Source: TraceSourceManual, TracedBy: "alice"}
	if err := tr.Validate(); err == nil {
		t.Error("expected error for zero version")
	}
}

func TestSecretTraceability_Validate_InvalidSource(t *testing.T) {
	tr := &SecretTraceability{Mount: "secret", Path: "app/db", Version: 1, Source: "bad", TracedBy: "alice"}
	if err := tr.Validate(); err == nil {
		t.Error("expected error for invalid source")
	}
}

func TestNewSecretTraceabilityRegistry_NotNil(t *testing.T) {
	if NewSecretTraceabilityRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestTraceabilityRegistry_Record_And_Get(t *testing.T) {
	reg := NewSecretTraceabilityRegistry()
	tr := &SecretTraceability{
		Mount:    "secret",
		Path:     "svc/key",
		Version:  2,
		Source:   TraceSourceImport,
		TracedBy: "ops",
	}
	if err := reg.Record(tr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("secret", "svc/key", 2)
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.TracedBy != "ops" {
		t.Errorf("unexpected traced_by: %s", got.TracedBy)
	}
}

func TestTraceabilityRegistry_Record_SetsTracedAt(t *testing.T) {
	reg := NewSecretTraceabilityRegistry()
	tr := &SecretTraceability{
		Mount: "secret", Path: "x", Version: 1,
		Source: TraceSourceGenerated, TracedBy: "gen",
	}
	_ = reg.Record(tr)
	if tr.TracedAt.IsZero() {
		t.Error("expected TracedAt to be set")
	}
}

func TestTraceabilityRegistry_Record_Invalid(t *testing.T) {
	reg := NewSecretTraceabilityRegistry()
	if err := reg.Record(&SecretTraceability{}); err == nil {
		t.Error("expected error for invalid record")
	}
}

func TestTraceabilityRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretTraceabilityRegistry()
	_, ok := reg.Get("secret", "missing", 1)
	if ok {
		t.Error("expected not found")
	}
}

func TestTraceabilityRegistry_Remove(t *testing.T) {
	reg := NewSecretTraceabilityRegistry()
	tr := &SecretTraceability{
		Mount: "secret", Path: "del", Version: 1,
		Source: TraceSourceMigrated, TracedBy: "admin",
		TracedAt: time.Now(),
	}
	_ = reg.Record(tr)
	reg.Remove("secret", "del", 1)
	if _, ok := reg.Get("secret", "del", 1); ok {
		t.Error("expected record to be removed")
	}
}
