package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTemp(t, "vault:\n  token: s.test\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected default address, got %q", cfg.Vault.Address)
	}
	if cfg.Output.Format != "text" {
		t.Errorf("expected default output format 'text', got %q", cfg.Output.Format)
	}
}

func TestLoad_FullConfig(t *testing.T) {
	yaml := `
vault:
  address: https://vault.example.com
  token: s.abc123
  namespace: prod
  tls_skip_verify: true
audit:
  enabled: true
  format: json
  path: /var/log/vaultdiff.log
output:
  format: json
  mask_secrets: true
`
	path := writeTemp(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Namespace != "prod" {
		t.Errorf("expected namespace 'prod', got %q", cfg.Vault.Namespace)
	}
	if !cfg.Output.MaskSecrets {
		t.Error("expected mask_secrets to be true")
	}
	if cfg.Audit.Path != "/var/log/vaultdiff.log" {
		t.Errorf("unexpected audit path: %q", cfg.Audit.Path)
	}
}

func TestLoad_InvalidOutputFormat(t *testing.T) {
	path := writeTemp(t, "output:\n  format: xml\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid output format")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTemp(t, ": : invalid: yaml:::")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
