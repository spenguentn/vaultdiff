package vault

import (
	"testing"
	"time"
)

func TestNewSecretDeprecationRegistry_NotNil(t *testing.T) {
	r := NewSecretDeprecationRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestDeprecationRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretDeprecationRegistry()
	d := SecretDeprecation{
		Mount:      "secret",
		Path:       "app/old",
		Status:     "deprecated",
		DeprecatedBy: "alice",
		SunsetAt:   time.Now().Add(72 * time.Hour),
	}
	if err := r.Set(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("secret", "app/old")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if got.Status != "deprecated" {
		t.Errorf("expected status 'deprecated', got %q", got.Status)
	}
}

func TestDeprecationRegistry_Set_SetsDeprecatedAt(t *testing.T) {
	r := NewSecretDeprecationRegistry()
	d := SecretDeprecation{
		Mount:      "secret",
		Path:       "app/old",
		Status:     "deprecated",
		DeprecatedBy: "alice",
		SunsetAt:   time.Now().Add(time.Hour),
	}
	_ = r.Set(d)
	got, _ := r.Get("secret", "app/old")
	if got.DeprecatedAt.IsZero() {
		t.Error("expected DeprecatedAt to be stamped")
	}
}

func TestDeprecationRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretDeprecationRegistry()
	if err := r.Set(SecretDeprecation{}); err == nil {
		t.Error("expected validation error")
	}
}

func TestDeprecationRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretDeprecationRegistry()
	_, ok := r.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestDeprecationRegistry_Remove(t *testing.T) {
	r := NewSecretDeprecationRegistry()
	d := SecretDeprecation{
		Mount:      "secret",
		Path:       "app/old",
		Status:     "deprecated",
		DeprecatedBy: "alice",
		SunsetAt:   time.Now().Add(time.Hour),
	}
	_ = r.Set(d)
	r.Remove("secret", "app/old")
	_, ok := r.Get("secret", "app/old")
	if ok {
		t.Error("expected record to be removed")
	}
}

func TestDeprecationRegistry_Sunsetted(t *testing.T) {
	r := NewSecretDeprecationRegistry()
	past := SecretDeprecation{
		Mount:      "secret",
		Path:       "app/past",
		Status:     "deprecated",
		DeprecatedBy: "alice",
		SunsetAt:   time.Now().Add(-time.Hour),
	}
	future := SecretDeprecation{
		Mount:      "secret",
		Path:       "app/future",
		Status:     "deprecated",
		DeprecatedBy: "alice",
		SunsetAt:   time.Now().Add(time.Hour),
	}
	_ = r.Set(past)
	_ = r.Set(future)
	sunset := r.Sunsetted()
	if len(sunset) != 1 {
		t.Errorf("expected 1 sunsetted record, got %d", len(sunset))
	}
	if sunset[0].Path != "app/past" {
		t.Errorf("unexpected path: %s", sunset[0].Path)
	}
}
