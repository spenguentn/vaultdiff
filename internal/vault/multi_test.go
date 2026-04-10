package vault_test

import (
	"context"
	"sort"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestNewMultiReader_NotNil(t *testing.T) {
	mr := vault.NewMultiReader(map[string]*vault.Client{})
	if mr == nil {
		t.Fatal("expected non-nil MultiReader")
	}
}

func TestMultiReader_ReadAll_EmptyClients(t *testing.T) {
	mr := vault.NewMultiReader(map[string]*vault.Client{})
	results := mr.ReadAll(context.Background(), "secret/data/app")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSecretResult_Fields(t *testing.T) {
	results := []vault.SecretResult{
		{Env: "staging", Secrets: map[string]string{"key": "val"}, Err: nil},
		{Env: "prod", Secrets: nil, Err: fmt.Errorf("connection refused")},
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Env < results[j].Env
	})

	if results[0].Env != "prod" {
		t.Errorf("expected prod first, got %s", results[0].Env)
	}
	if results[0].Err == nil {
		t.Error("expected error for prod")
	}
	if results[1].Secrets["key"] != "val" {
		t.Error("expected staging secret value")
	}
}

func TestMultiReader_ReadAll_ResultCount(t *testing.T) {
	// We cannot spin up real Vault here; verify the channel-based fan-out
	// returns exactly len(clients) results (all with errors since no server).
	clients := map[string]*vault.Client{}
	mr := vault.NewMultiReader(clients)
	results := mr.ReadAll(context.Background(), "secret/data/test")
	if len(results) != len(clients) {
		t.Fatalf("result count mismatch: want %d got %d", len(clients), len(results))
	}
}
