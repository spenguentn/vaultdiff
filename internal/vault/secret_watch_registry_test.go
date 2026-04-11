package vault

import (
	"testing"
	"time"
)

func TestNewSecretWatchRegistry_NotNil(t *testing.T) {
	r := NewSecretWatchRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestWatchRegistry_Register_And_Get(t *testing.T) {
	r := NewSecretWatchRegistry()
	cfg := validWatchConfig()
	if err := r.Register(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get(cfg.Mount, cfg.Path)
	if !ok {
		t.Fatal("expected watch config to be found")
	}
	if got.Path != cfg.Path {
		t.Errorf("path mismatch: got %s, want %s", got.Path, cfg.Path)
	}
}

func TestWatchRegistry_Register_Invalid(t *testing.T) {
	r := NewSecretWatchRegistry()
	cfg := validWatchConfig()
	cfg.Mount = ""
	if err := r.Register(cfg); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestWatchRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretWatchRegistry()
	_, ok := r.Get("secret", "missing/path")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestWatchRegistry_Remove(t *testing.T) {
	r := NewSecretWatchRegistry()
	cfg := validWatchConfig()
	_ = r.Register(cfg)
	r.Remove(cfg.Mount, cfg.Path)
	_, ok := r.Get(cfg.Mount, cfg.Path)
	if ok {
		t.Fatal("expected watch to be removed")
	}
}

func TestWatchRegistry_List(t *testing.T) {
	r := NewSecretWatchRegistry()
	cfg1 := validWatchConfig()
	cfg2 := SecretWatchConfig{Mount: "kv", Path: "db/creds", Interval: 5 * time.Second, OnChange: noop}
	_ = r.Register(cfg1)
	_ = r.Register(cfg2)
	list := r.List()
	if len(list) != 2 {
		t.Errorf("expected 2 watches, got %d", len(list))
	}
}
