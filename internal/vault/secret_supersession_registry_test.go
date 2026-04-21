package vault

import (
	"testing"
	"time"
)

func TestNewSecretSupersessionRegistry_NotNil(t *testing.T) {
	r := NewSecretSupersessionRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestSupersessionRegistry_Record_And_Get(t *testing.T) {
	r := NewSecretSupersessionRegistry()
	s := SecretSupersession{
		Mount:          "secret",
		Path:           "app/db",
		SupersededBy:   "app/db-v2",
		Reason:         SupersessionReasonRotated,
		SupersededAt:   time.Now().UTC(),
	}
	if err := r.Record(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/db")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.SupersededBy != "app/db-v2" {
		t.Errorf("expected SupersededBy app/db-v2, got %s", got.SupersededBy)
	}
}

func TestSupersessionRegistry_Record_SetsSupersededAt(t *testing.T) {
	r := NewSecretSupersessionRegistry()
	s := SecretSupersession{
		Mount:        "secret",
		Path:         "app/key",
		SupersededBy: "app/key-v2",
		Reason:       SupersessionReasonRotated,
	}
	if err := r.Record(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := r.Get("secret", "app/key")
	if got.SupersededAt.IsZero() {
		t.Error("expected SupersededAt to be set automatically")
	}
}

func TestSupersessionRegistry_Record_Invalid(t *testing.T) {
	r := NewSecretSupersessionRegistry()
	s := SecretSupersession{} // missing required fields
	if err := r.Record(s); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestSupersessionRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretSupersessionRegistry()
	_, ok := r.Get("secret", "nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestSupersessionRegistry_Remove(t *testing.T) {
	r := NewSecretSupersessionRegistry()
	s := SecretSupersession{
		Mount:        "secret",
		Path:         "app/cred",
		SupersededBy: "app/cred-v2",
		Reason:       SupersessionReasonDeprecated,
		SupersededAt: time.Now().UTC(),
	}
	_ = r.Record(s)
	r.Remove("secret", "app/cred")
	_, ok := r.Get("secret", "app/cred")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestSupersessionRegistry_All(t *testing.T) {
	r := NewSecretSupersessionRegistry()
	for _, path := range []string{"a", "b", "c"} {
		_ = r.Record(SecretSupersession{
			Mount:        "secret",
			Path:         path,
			SupersededBy: path + "-v2",
			Reason:       SupersessionReasonRotated,
			SupersededAt: time.Now().UTC(),
		})
	}
	if len(r.All()) != 3 {
		t.Errorf("expected 3 entries, got %d", len(r.All()))
	}
}
