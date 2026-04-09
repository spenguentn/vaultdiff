package vault_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var validLeft = vault.Environment{
	Name:      "staging",
	Address:   "https://vault.staging.example.com",
	MountPath: "secret",
	Token:     "s.staging",
}

var validRight = vault.Environment{
	Name:      "production",
	Address:   "https://vault.prod.example.com",
	MountPath: "secret",
	Token:     "s.production",
}

func TestNewEnvPair_Valid(t *testing.T) {
	pair, err := vault.NewEnvPair(validLeft, validRight)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pair.Left.Name != "staging" {
		t.Errorf("expected left name staging, got %s", pair.Left.Name)
	}
	if pair.Right.Name != "production" {
		t.Errorf("expected right name production, got %s", pair.Right.Name)
	}
}

func TestNewEnvPair_InvalidLeft(t *testing.T) {
	bad := vault.Environment{Name: "", Address: "https://vault.example.com", MountPath: "secret", Token: "t"}
	_, err := vault.NewEnvPair(bad, validRight)
	if err == nil {
		t.Fatal("expected error for invalid left environment")
	}
}

func TestNewEnvPair_InvalidRight(t *testing.T) {
	bad := vault.Environment{Name: "prod", Address: "", MountPath: "secret", Token: "t"}
	_, err := vault.NewEnvPair(validLeft, bad)
	if err == nil {
		t.Fatal("expected error for invalid right environment")
	}
}

func TestEnvPair_Names(t *testing.T) {
	pair, _ := vault.NewEnvPair(validLeft, validRight)
	got := pair.Names()
	want := "staging → production"
	if got != want {
		t.Errorf("Names() = %q, want %q", got, want)
	}
}

func TestEnvPair_SameMount_True(t *testing.T) {
	pair, _ := vault.NewEnvPair(validLeft, validRight)
	if !pair.SameMount() {
		t.Error("expected SameMount to be true")
	}
}

func TestEnvPair_SameMount_False(t *testing.T) {
	diffMount := validRight
	diffMount.MountPath = "kv"
	pair, _ := vault.NewEnvPair(validLeft, diffMount)
	if pair.SameMount() {
		t.Error("expected SameMount to be false")
	}
}
