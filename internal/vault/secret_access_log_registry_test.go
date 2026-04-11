package vault

import (
	"testing"
	"time"
)

func TestNewSecretAccessLogRegistry_NotNil(t *testing.T) {
	r := NewSecretAccessLogRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestRegistry_Record_And_Get(t *testing.T) {
	r := NewSecretAccessLogRegistry()
	e := SecretAccessEntry{
		Mount: "secret", Path: "app/key",
		EventType: AccessEventWrite, Actor: "alice",
		Timestamp: time.Now().UTC(),
	}
	if err := r.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, ok := r.Get("secret", "app/key")
	if !ok || len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestRegistry_Record_SetsTimestamp(t *testing.T) {
	r := NewSecretAccessLogRegistry()
	e := SecretAccessEntry{
		Mount: "secret", Path: "app/key",
		EventType: AccessEventRead, Actor: "bob",
	}
	_ = r.Record(e)
	entries, _ := r.Get("secret", "app/key")
	if entries[0].Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}
}

func TestRegistry_Record_InvalidEntry(t *testing.T) {
	r := NewSecretAccessLogRegistry()
	e := SecretAccessEntry{Path: "app/key", EventType: AccessEventRead, Actor: "bob"}
	if err := r.Record(e); err == nil {
		t.Error("expected error for invalid entry")
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretAccessLogRegistry()
	_, ok := r.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestRegistry_All_ReturnsAll(t *testing.T) {
	r := NewSecretAccessLogRegistry()
	for _, path := range []string{"a", "b", "c"} {
		_ = r.Record(SecretAccessEntry{
			Mount: "kv", Path: path,
			EventType: AccessEventList, Actor: "svc",
			Timestamp: time.Now().UTC(),
		})
	}
	if got := len(r.All()); got != 3 {
		t.Errorf("expected 3 entries, got %d", got)
	}
}

func TestRegistry_Clear(t *testing.T) {
	r := NewSecretAccessLogRegistry()
	_ = r.Record(SecretAccessEntry{
		Mount: "kv", Path: "x",
		EventType: AccessEventDelete, Actor: "admin",
		Timestamp: time.Now().UTC(),
	})
	r.Clear()
	if got := len(r.All()); got != 0 {
		t.Errorf("expected 0 after clear, got %d", got)
	}
}
