package vault

import (
	"testing"
)

func baseEnv() *Environment {
	return &Environment{
		Name:      "staging",
		Address:   "https://vault.staging.example.com",
		MountPath: "secret",
		Token:     "s.abc123",
	}
}

func TestEnvironment_Validate_Valid(t *testing.T) {
	env := baseEnv()
	if err := env.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestEnvironment_Validate_MissingName(t *testing.T) {
	env := baseEnv()
	env.Name = ""
	if err := env.Validate(); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestEnvironment_Validate_MissingAddress(t *testing.T) {
	env := baseEnv()
	env.Address = ""
	if err := env.Validate(); err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestEnvironment_Validate_MissingMountPath(t *testing.T) {
	env := baseEnv()
	env.MountPath = ""
	if err := env.Validate(); err == nil {
		t.Fatal("expected error for missing mount path")
	}
}

func TestEnvironment_Validate_MissingToken(t *testing.T) {
	env := baseEnv()
	env.Token = ""
	if err := env.Validate(); err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestEnvironment_SecretPath(t *testing.T) {
	env := baseEnv()
	got := env.SecretPath("myapp/db")
	want := "secret/data/myapp/db"
	if got != want {
		t.Errorf("SecretPath() = %q, want %q", got, want)
	}
}

func TestEnvironment_MetadataPath(t *testing.T) {
	env := baseEnv()
	got := env.MetadataPath("myapp/db")
	want := "secret/metadata/myapp/db"
	if got != want {
		t.Errorf("MetadataPath() = %q, want %q", got, want)
	}
}
