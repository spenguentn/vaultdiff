package vault_test

import (
	"fmt"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func validEnv(name string) vault.Environment {
	return vault.Environment{
		Name:      name,
		Address:   "http://127.0.0.1:8200",
		MountPath: "secret",
	}
}

func TestNewMultiClientBuilder_NotNil(t *testing.T) {
	b := vault.NewMultiClientBuilder(nil)
	if b == nil {
		t.Fatal("expected non-nil builder")
	}
}

func TestMultiClientBuilder_Build_Empty(t *testing.T) {
	clients, err := vault.NewMultiClientBuilder([]vault.Environment{}).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(clients) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(clients))
	}
}

func TestMultiClientBuilder_Build_InvalidEnv(t *testing.T) {
	envs := []vault.Environment{
		{Name: "", Address: "http://127.0.0.1:8200", MountPath: "secret"},
	}
	_, err := vault.NewMultiClientBuilder(envs).Build()
	if err == nil {
		t.Fatal("expected error for invalid environment")
	}
}

func TestMultiClientBuilder_WithToken_AppliedPerEnv(t *testing.T) {
	envs := []vault.Environment{validEnv("dev"), validEnv("prod")}
	b := vault.NewMultiClientBuilder(envs).
		WithToken("dev", "dev-token").
		WithToken("prod", "prod-token")

	clients, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(clients) != 2 {
		t.Fatalf("expected 2 clients, got %d", len(clients))
	}
	for _, name := range []string{"dev", "prod"} {
		if _, ok := clients[name]; !ok {
			t.Errorf("missing client for env %q", name)
		}
	}
	_ = fmt.Sprintf("clients built: %v", clients)
}

func TestMultiClientBuilder_Build_DuplicateEnvNames(t *testing.T) {
	envs := []vault.Environment{validEnv("staging"), validEnv("staging")}
	clients, err := vault.NewMultiClientBuilder(envs).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// second entry overwrites first; map still has 1 key
	if len(clients) != 1 {
		t.Fatalf("expected 1 client (deduped), got %d", len(clients))
	}
}
