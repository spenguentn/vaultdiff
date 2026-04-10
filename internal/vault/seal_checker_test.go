package vault

import (
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestNewSealChecker_Valid(t *testing.T) {
	client, _ := vaultapi.NewClient(vaultapi.DefaultConfig())
	sc, err := NewSealChecker(client, 3*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc == nil {
		t.Fatal("expected non-nil SealChecker")
	}
}

func TestNewSealChecker_NilClient(t *testing.T) {
	_, err := NewSealChecker(nil, 3*time.Second)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestNewSealChecker_DefaultTimeout(t *testing.T) {
	client, _ := vaultapi.NewClient(vaultapi.DefaultConfig())
	sc, err := NewSealChecker(client, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", sc.timeout)
	}
}

func TestNewSealChecker_NegativeTimeout(t *testing.T) {
	client, _ := vaultapi.NewClient(vaultapi.DefaultConfig())
	sc, err := NewSealChecker(client, -1*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", sc.timeout)
	}
}
