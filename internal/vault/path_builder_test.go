package vault

import (
	"testing"
)

func baseBuilder() *PathBuilder {
	env := Environment{
		Name:      "staging",
		Address:   "https://vault.example.com",
		MountPath: "secret",
		Token:     "tok",
	}
	return NewPathBuilder(env)
}

func TestPathBuilder_Secret_Valid(t *testing.T) {
	b := baseBuilder()
	sp, err := b.Secret("app/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sp.DataPath() != "secret/data/app/config" {
		t.Errorf("expected secret/data/app/config, got %s", sp.DataPath())
	}
}

func TestPathBuilder_Secret_TrimsSlashes(t *testing.T) {
	b := baseBuilder()
	sp, err := b.Secret("/app/config/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sp.DataPath() != "secret/data/app/config" {
		t.Errorf("expected trimmed path, got %s", sp.DataPath())
	}
}

func TestPathBuilder_Secret_EmptyPath(t *testing.T) {
	b := baseBuilder()
	_, err := b.Secret("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestPathBuilder_MustSecret_Panics(t *testing.T) {
	b := baseBuilder()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty path")
		}
	}()
	b.MustSecret("")
}

func TestPathBuilder_Batch_Valid(t *testing.T) {
	b := baseBuilder()
	paths := []string{"app/db", "app/cache", "infra/tls"}
	result, err := b.Batch(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(paths) {
		t.Errorf("expected %d paths, got %d", len(paths), len(result))
	}
}

func TestPathBuilder_Batch_ErrorOnEmpty(t *testing.T) {
	b := baseBuilder()
	_, err := b.Batch([]string{"app/db", "", "infra/tls"})
	if err == nil {
		t.Fatal("expected error for empty path in batch")
	}
}

func TestPathBuilder_EnvPrefix(t *testing.T) {
	b := baseBuilder()
	if b.EnvPrefix() != "[staging]" {
		t.Errorf("expected [staging], got %s", b.EnvPrefix())
	}
}
