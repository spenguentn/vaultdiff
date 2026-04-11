package vault

import (
	"testing"
	"time"
)

func noop(_ SecretWatchEvent) {}

func validWatchConfig() SecretWatchConfig {
	return SecretWatchConfig{
		Mount:    "secret",
		Path:     "app/config",
		Interval: 10 * time.Second,
		OnChange: noop,
	}
}

func TestSecretWatchConfig_Validate_Valid(t *testing.T) {
	if err := validWatchConfig().Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSecretWatchConfig_Validate_MissingMount(t *testing.T) {
	cfg := validWatchConfig()
	cfg.Mount = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestSecretWatchConfig_Validate_MissingPath(t *testing.T) {
	cfg := validWatchConfig()
	cfg.Path = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestSecretWatchConfig_Validate_ZeroInterval(t *testing.T) {
	cfg := validWatchConfig()
	cfg.Interval = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestSecretWatchConfig_Validate_NilHandler(t *testing.T) {
	cfg := validWatchConfig()
	cfg.OnChange = nil
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for nil OnChange")
	}
}

func TestSecretWatchEvent_IsChanged_True(t *testing.T) {
	e := SecretWatchEvent{PrevHash: "abc", CurrHash: "xyz"}
	if !e.IsChanged() {
		t.Fatal("expected IsChanged to return true")
	}
}

func TestSecretWatchEvent_IsChanged_False(t *testing.T) {
	e := SecretWatchEvent{PrevHash: "abc", CurrHash: "abc"}
	if e.IsChanged() {
		t.Fatal("expected IsChanged to return false")
	}
}
