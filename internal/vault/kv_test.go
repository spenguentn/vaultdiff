package vault

import (
	"testing"
)

func baseKVConfig() KVConfig {
	return KVConfig{
		MountPath: "secret",
		Version:   KVv2,
	}
}

func TestKVConfig_Validate_Valid(t *testing.T) {
	if err := baseKVConfig().Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestKVConfig_Validate_MissingMountPath(t *testing.T) {
	cfg := baseKVConfig()
	cfg.MountPath = ""
	if err := cfg.Validate(); err != ErrMissingMountPath {
		t.Fatalf("expected ErrMissingMountPath, got %v", err)
	}
}

func TestKVConfig_Validate_InvalidVersion(t *testing.T) {
	cfg := baseKVConfig()
	cfg.Version = KVVersion(99)
	if err := cfg.Validate(); err != ErrInvalidKVVersion {
		t.Fatalf("expected ErrInvalidKVVersion, got %v", err)
	}
}

func TestKVConfig_IsVersioned_V2(t *testing.T) {
	if !baseKVConfig().IsVersioned() {
		t.Fatal("expected KVv2 to be versioned")
	}
}

func TestKVConfig_IsVersioned_V1(t *testing.T) {
	cfg := KVConfig{MountPath: "kv", Version: KVv1}
	if cfg.IsVersioned() {
		t.Fatal("expected KVv1 to not be versioned")
	}
}

func TestKVConfig_DataPrefix_V2(t *testing.T) {
	got := baseKVConfig().DataPrefix()
	want := "secret/data"
	if got != want {
		t.Fatalf("DataPrefix: got %q, want %q", got, want)
	}
}

func TestKVConfig_DataPrefix_V1(t *testing.T) {
	cfg := KVConfig{MountPath: "kv", Version: KVv1}
	got := cfg.DataPrefix()
	if got != "kv" {
		t.Fatalf("DataPrefix V1: got %q, want %q", got, "kv")
	}
}

func TestKVConfig_MetadataPrefix_V2(t *testing.T) {
	got := baseKVConfig().MetadataPrefix()
	want := "secret/metadata"
	if got != want {
		t.Fatalf("MetadataPrefix: got %q, want %q", got, want)
	}
}

func TestKVConfig_MetadataPrefix_V1(t *testing.T) {
	cfg := KVConfig{MountPath: "kv", Version: KVv1}
	if got := cfg.MetadataPrefix(); got != "" {
		t.Fatalf("MetadataPrefix V1: expected empty string, got %q", got)
	}
}
