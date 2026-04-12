package vault

import (
	"testing"
)

func TestDefaultRedactConfig_MaskString(t *testing.T) {
	cfg := DefaultRedactConfig()
	if cfg.MaskString != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", cfg.MaskString)
	}
}

func TestShouldRedact_MatchesPassword(t *testing.T) {
	cfg := DefaultRedactConfig()
	if !cfg.ShouldRedact("password") {
		t.Error("expected password to be redacted")
	}
}

func TestShouldRedact_MatchesToken(t *testing.T) {
	cfg := DefaultRedactConfig()
	if !cfg.ShouldRedact("auth_token") {
		t.Error("expected auth_token to be redacted")
	}
}

func TestShouldRedact_DoesNotMatchPlainKey(t *testing.T) {
	cfg := DefaultRedactConfig()
	if cfg.ShouldRedact("database_host") {
		t.Error("expected database_host not to be redacted")
	}
}

func TestShouldRedact_MatchesKeyPrefix(t *testing.T) {
	cfg := DefaultRedactConfig()
	cfg.KeyPrefixes = []string{"internal_"}
	if !cfg.ShouldRedact("internal_config") {
		t.Error("expected internal_config to be redacted via prefix")
	}
}

func TestApply_MaskMode(t *testing.T) {
	cfg := DefaultRedactConfig()
	data := map[string]string{
		"password": "s3cr3t",
		"host":     "localhost",
	}
	out := cfg.Apply(data)
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out["password"])
	}
	if out["host"] != "localhost" {
		t.Errorf("expected localhost, got %q", out["host"])
	}
}

func TestApply_RemoveMode(t *testing.T) {
	cfg := DefaultRedactConfig()
	cfg.Mode = RedactRemove
	data := map[string]string{
		"api_key": "abc123",
		"region":  "us-east-1",
	}
	out := cfg.Apply(data)
	if _, ok := out["api_key"]; ok {
		t.Error("expected api_key to be removed")
	}
	if out["region"] != "us-east-1" {
		t.Errorf("expected region to be preserved")
	}
}

func TestApply_PartialMode(t *testing.T) {
	cfg := DefaultRedactConfig()
	cfg.Mode = RedactPartial
	data := map[string]string{
		"password": "abcdef",
	}
	out := cfg.Apply(data)
	if out["password"] != "a****f" {
		t.Errorf("expected a****f, got %q", out["password"])
	}
}

func TestPartialMask_ShortValue(t *testing.T) {
	result := partialMask("ab")
	if result != "**" {
		t.Errorf("expected **, got %q", result)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	cfg := DefaultRedactConfig()
	original := map[string]string{"secret": "val"}
	_ = cfg.Apply(original)
	if original["secret"] != "val" {
		t.Error("original map should not be mutated")
	}
}
