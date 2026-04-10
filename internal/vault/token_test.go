package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveToken_Direct(t *testing.T) {
	tok, err := ResolveToken("s.direct123", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value != "s.direct123" {
		t.Errorf("expected s.direct123, got %q", tok.Value)
	}
	if tok.Source != TokenSourceDirect {
		t.Errorf("expected TokenSourceDirect, got %v", tok.Source)
	}
}

func TestResolveToken_EnvVar(t *testing.T) {
	t.Setenv("VAULT_TOKEN_TEST", "s.envtoken")
	tok, err := ResolveToken("", "VAULT_TOKEN_TEST", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value != "s.envtoken" {
		t.Errorf("expected s.envtoken, got %q", tok.Value)
	}
	if tok.Source != TokenSourceEnv {
		t.Errorf("expected TokenSourceEnv, got %v", tok.Source)
	}
}

func TestResolveToken_File(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "token")
	if err := os.WriteFile(p, []byte("s.filetoken\n"), 0600); err != nil {
		t.Fatal(err)
	}
	tok, err := ResolveToken("", "", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value != "s.filetoken" {
		t.Errorf("expected s.filetoken, got %q", tok.Value)
	}
	if tok.Source != TokenSourceFile {
		t.Errorf("expected TokenSourceFile, got %v", tok.Source)
	}
}

func TestResolveToken_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "token")
	if err := os.WriteFile(p, []byte("   \n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := ResolveToken("", "", p)
	if err == nil {
		t.Error("expected error for empty token file")
	}
}

func TestResolveToken_NoneProvided(t *testing.T) {
	_, err := ResolveToken("", "", "")
	if err == nil {
		t.Error("expected error when no token source is provided")
	}
}

func TestResolveToken_DirectTakesPriority(t *testing.T) {
	t.Setenv("VAULT_TOKEN_PRI", "s.shouldbeignored")
	tok, err := ResolveToken("s.direct", "VAULT_TOKEN_PRI", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Source != TokenSourceDirect {
		t.Errorf("expected TokenSourceDirect, got %v", tok.Source)
	}
}
