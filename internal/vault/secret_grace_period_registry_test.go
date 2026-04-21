package vault

import (
	"testing"
	"time"
)

func TestNewSecretGracePeriodRegistry_NotNil(t *testing.T) {
	r := NewSecretGracePeriodRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestGracePeriodRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretGracePeriodRegistry()
	g := SecretGracePeriod{
		Mount:     "secret",
		Path:      "app/key",
		Status:    "active",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := r.Set(g); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/key")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if got.Status != "active" {
		t.Errorf("expected status 'active', got %q", got.Status)
	}
}

func TestGracePeriodRegistry_Set_SetsCreatedAt(t *testing.T) {
	r := NewSecretGracePeriodRegistry()
	g := SecretGracePeriod{
		Mount:     "secret",
		Path:      "app/key",
		Status:    "active",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	_ = r.Set(g)
	got, _ := r.Get("secret", "app/key")
	if got.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestGracePeriodRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretGracePeriodRegistry()
	g := SecretGracePeriod{} // missing required fields
	if err := r.Set(g); err == nil {
		t.Error("expected validation error")
	}
}

func TestGracePeriodRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretGracePeriodRegistry()
	_, ok := r.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestGracePeriodRegistry_Remove(t *testing.T) {
	r := NewSecretGracePeriodRegistry()
	g := SecretGracePeriod{
		Mount:     "secret",
		Path:      "app/key",
		Status:    "active",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	_ = r.Set(g)
	r.Remove("secret", "app/key")
	_, ok := r.Get("secret", "app/key")
	if ok {
		t.Error("expected record to be removed")
	}
}

func TestGracePeriodRegistry_Expired(t *testing.T) {
	r := NewSecretGracePeriodRegistry()
	past := SecretGracePeriod{
		Mount:     "secret",
		Path:      "old/key",
		Status:    "active",
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	future := SecretGracePeriod{
		Mount:     "secret",
		Path:      "new/key",
		Status:    "active",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	_ = r.Set(past)
	_ = r.Set(future)
	expired := r.Expired()
	if len(expired) != 1 {
		t.Errorf("expected 1 expired record, got %d", len(expired))
	}
	if expired[0].Path != "old/key" {
		t.Errorf("unexpected expired path: %s", expired[0].Path)
	}
}
