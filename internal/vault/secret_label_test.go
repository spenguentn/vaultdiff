package vault

import (
	"testing"
)

func TestSecretLabel_FullPath(t *testing.T) {
	l := SecretLabel{Mount: "secret", Path: "app/db"}
	if got := l.FullPath(); got != "secret/app/db" {
		t.Fatalf("expected secret/app/db, got %s", got)
	}
}

func TestSecretLabel_Validate_Valid(t *testing.T) {
	l := SecretLabel{Mount: "secret", Path: "app/db", Key: "env", Value: "prod", CreatedBy: "alice"}
	if err := l.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretLabel_Validate_MissingMount(t *testing.T) {
	l := SecretLabel{Path: "app/db", Key: "env", CreatedBy: "alice"}
	if err := l.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestSecretLabel_Validate_MissingPath(t *testing.T) {
	l := SecretLabel{Mount: "secret", Key: "env", CreatedBy: "alice"}
	if err := l.Validate(); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestSecretLabel_Validate_MissingKey(t *testing.T) {
	l := SecretLabel{Mount: "secret", Path: "app/db", CreatedBy: "alice"}
	if err := l.Validate(); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestSecretLabel_Validate_MissingCreatedBy(t *testing.T) {
	l := SecretLabel{Mount: "secret", Path: "app/db", Key: "env"}
	if err := l.Validate(); err == nil {
		t.Fatal("expected error for missing created_by")
	}
}
