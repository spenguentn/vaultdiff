package vault

import (
	"testing"
	"time"
)

func TestNewSecretPinRegistry_NotNil(t *testing.T) {
	r := NewSecretPinRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestPinRegistry_Pin_And_Get(t *testing.T) {
	r := NewSecretPinRegistry()
	if err := r.Pin(basePin); err != nil {
		t.Fatalf("Pin() error: %v", err)
	}
	got, ok := r.Get(basePin.Mount, basePin.Path)
	if !ok {
		t.Fatal("expected pin to be found")
	}
	if got.Version != basePin.Version {
		t.Errorf("Version = %d, want %d", got.Version, basePin.Version)
	}
}

func TestPinRegistry_Pin_SetsPinnedAt(t *testing.T) {
	r := NewSecretPinRegistry()
	p := basePin
	p.PinnedAt = time.Time{}
	_ = r.Pin(p)
	got, _ := r.Get(p.Mount, p.Path)
	if got.PinnedAt.IsZero() {
		t.Error("expected PinnedAt to be set automatically")
	}
}

func TestPinRegistry_Pin_Invalid(t *testing.T) {
	r := NewSecretPinRegistry()
	p := basePin
	p.Mount = ""
	if err := r.Pin(p); err == nil {
		t.Error("expected error for invalid pin")
	}
}

func TestPinRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretPinRegistry()
	_, ok := r.Get("secret", "missing/path")
	if ok {
		t.Error("expected not found")
	}
}

func TestPinRegistry_Unpin(t *testing.T) {
	r := NewSecretPinRegistry()
	_ = r.Pin(basePin)
	removed := r.Unpin(basePin.Mount, basePin.Path)
	if !removed {
		t.Error("expected Unpin to return true")
	}
	_, ok := r.Get(basePin.Mount, basePin.Path)
	if ok {
		t.Error("expected pin to be removed")
	}
}

func TestPinRegistry_Unpin_NotFound(t *testing.T) {
	r := NewSecretPinRegistry()
	if r.Unpin("secret", "ghost/path") {
		t.Error("expected false for non-existent pin")
	}
}

func TestPinRegistry_IsPinned_Active(t *testing.T) {
	r := NewSecretPinRegistry()
	_ = r.Pin(basePin)
	if !r.IsPinned(basePin.Mount, basePin.Path) {
		t.Error("expected path to be pinned")
	}
}

func TestPinRegistry_IsPinned_Expired(t *testing.T) {
	r := NewSecretPinRegistry()
	p := basePin
	p.ExpiresAt = time.Now().Add(-time.Minute)
	_ = r.Pin(p)
	if r.IsPinned(p.Mount, p.Path) {
		t.Error("expected expired pin to report not pinned")
	}
}

func TestPinRegistry_All(t *testing.T) {
	r := NewSecretPinRegistry()
	_ = r.Pin(basePin)
	if len(r.All()) != 1 {
		t.Errorf("All() len = %d, want 1", len(r.All()))
	}
}
