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
	tests := []struct {
		name string
		path string
		want string
	}{
		{"simple path", "myapp/db", "secret/data/myapp/db"},
		{"nested path", "myapp/prod/db", "secret/data/myapp/prod/db"},
		{"single segment", "myapp", "secret/data/myapp"},
	}
	env := baseEnv()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := env.SecretPath(tt.path); got != tt.want {
				t.Errorf("SecretPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
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
