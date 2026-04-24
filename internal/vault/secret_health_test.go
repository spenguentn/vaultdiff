package vault

import (
	"testing"
	"time"
)

func TestIsValidSecretHealthStatus_Known(t *testing.T) {
	for _, s := range []SecretHealthStatus{
		SecretHealthOK, SecretHealthDegraded, SecretHealthCritical, SecretHealthUnknown,
	} {
		if !IsValidSecretHealthStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidSecretHealthStatus_Unknown(t *testing.T) {
	if IsValidSecretHealthStatus("bogus") {
		t.Error("expected 'bogus' to be invalid")
	}
}

func TestSecretHealth_FullPath(t *testing.T) {
	h := SecretHealth{Mount: "secret", Path: "app/db"}
	if got := h.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretHealth_IsHealthy_True(t *testing.T) {
	h := SecretHealth{Status: SecretHealthOK}
	if !h.IsHealthy() {
		t.Error("expected IsHealthy to return true")
	}
}

func TestSecretHealth_IsHealthy_False(t *testing.T) {
	h := SecretHealth{Status: SecretHealthCritical}
	if h.IsHealthy() {
		t.Error("expected IsHealthy to return false")
	}
}

func TestSecretHealth_Validate_Valid(t *testing.T) {
	h := SecretHealth{
		Mount:     "secret",
		Path:      "app/key",
		Status:    SecretHealthOK,
		CheckedBy: "ci-bot",
		CheckedAt: time.Now(),
	}
	if err := h.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretHealth_Validate_MissingMount(t *testing.T) {
	h := SecretHealth{Path: "app/key", Status: SecretHealthOK, CheckedBy: "bot", CheckedAt: time.Now()}
	if err := h.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretHealth_Validate_InvalidStatus(t *testing.T) {
	h := SecretHealth{Mount: "m", Path: "p", Status: "nope", CheckedBy: "bot", CheckedAt: time.Now()}
	if err := h.Validate(); err == nil {
		t.Error("expected error for invalid status")
	}
}

func TestNewSecretHealthRegistry_NotNil(t *testing.T) {
	if NewSecretHealthRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestHealthRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretHealthRegistry()
	h := SecretHealth{Mount: "secret", Path: "app/key", Status: SecretHealthOK, CheckedBy: "bot", CheckedAt: time.Now()}
	if err := r.Set(h); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/key")
	if !ok {
		t.Fatal("expected record to be found")
	}
	if got.Status != SecretHealthOK {
		t.Errorf("unexpected status: %s", got.Status)
	}
}

func TestHealthRegistry_Set_SetsCheckedAt(t *testing.T) {
	r := NewSecretHealthRegistry()
	h := SecretHealth{Mount: "m", Path: "p", Status: SecretHealthDegraded, CheckedBy: "bot"}
	if err := r.Set(h); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := r.Get("m", "p")
	if got.CheckedAt.IsZero() {
		t.Error("expected CheckedAt to be stamped")
	}
}

func TestHealthRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretHealthRegistry()
	_, ok := r.Get("x", "y")
	if ok {
		t.Error("expected not found")
	}
}

func TestHealthRegistry_Remove(t *testing.T) {
	r := NewSecretHealthRegistry()
	h := SecretHealth{Mount: "m", Path: "p", Status: SecretHealthOK, CheckedBy: "bot", CheckedAt: time.Now()}
	_ = r.Set(h)
	r.Remove("m", "p")
	_, ok := r.Get("m", "p")
	if ok {
		t.Error("expected record to be removed")
	}
}
