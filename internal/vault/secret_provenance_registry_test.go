package vault

import (
	"testing"
	"time"
)

func TestNewSecretProvenanceRegistry_NotNil(t *testing.T) {
	r := NewSecretProvenanceRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestProvenanceRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretProvenanceRegistry()
	p := SecretProvenance{
		Mount:  "secret",
		Path:   "app/key",
		Source: ProvenanceSourceManual,
		Author: "alice",
	}
	if err := r.Set(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/key")
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.Author != "alice" {
		t.Errorf("expected author alice, got %s", got.Author)
	}
}

func TestProvenanceRegistry_Set_SetsRecordedAt(t *testing.T) {
	r := NewSecretProvenanceRegistry()
	p := SecretProvenance{
		Mount:  "secret",
		Path:   "app/key",
		Source: ProvenanceSourceManual,
		Author: "bob",
	}
	before := time.Now().UTC()
	_ = r.Set(p)
	got, _ := r.Get("secret", "app/key")
	if got.RecordedAt.Before(before) {
		t.Error("expected RecordedAt to be set to current time")
	}
}

func TestProvenanceRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretProvenanceRegistry()
	p := SecretProvenance{} // missing required fields
	if err := r.Set(p); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestProvenanceRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretProvenanceRegistry()
	_, ok := r.Get("secret", "missing/key")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestProvenanceRegistry_Remove(t *testing.T) {
	r := NewSecretProvenanceRegistry()
	p := SecretProvenance{
		Mount:  "secret",
		Path:   "app/key",
		Source: ProvenanceSourceManual,
		Author: "carol",
	}
	_ = r.Set(p)
	r.Remove("secret", "app/key")
	_, ok := r.Get("secret", "app/key")
	if ok {
		t.Fatal("expected record to be removed")
	}
}

func TestProvenanceRegistry_All(t *testing.T) {
	r := NewSecretProvenanceRegistry()
	for _, path := range []string{"a", "b", "c"} {
		_ = r.Set(SecretProvenance{
			Mount:  "secret",
			Path:   path,
			Source: ProvenanceSourceManual,
			Author: "dave",
		})
	}
	if len(r.All()) != 3 {
		t.Errorf("expected 3 records, got %d", len(r.All()))
	}
}
