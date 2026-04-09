package snapshot_test

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/snapshot"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
	}
}

func TestNew_FieldsSet(t *testing.T) {
	s := snapshot.New("secret/myapp", 3, baseSecrets(), snapshot.Meta{Environment: "prod"})
	if s.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", s.Path)
	}
	if s.Version != 3 {
		t.Errorf("expected version 3, got %d", s.Version)
	}
	if s.Meta.Environment != "prod" {
		t.Errorf("expected env prod, got %s", s.Meta.Environment)
	}
	if s.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestNew_NilSecretsDefaults(t *testing.T) {
	s := snapshot.New("secret/app", 1, nil, snapshot.Meta{})
	if s.Secrets == nil {
		t.Error("expected non-nil secrets map")
	}
}

func TestKeyCount(t *testing.T) {
	s := snapshot.New("secret/app", 1, baseSecrets(), snapshot.Meta{})
	if s.KeyCount() != 2 {
		t.Errorf("expected 2 keys, got %d", s.KeyCount())
	}
}

func TestHasKey(t *testing.T) {
	s := snapshot.New("secret/app", 1, baseSecrets(), snapshot.Meta{})
	if !s.HasKey("DB_HOST") {
		t.Error("expected HasKey to return true for DB_HOST")
	}
	if s.HasKey("MISSING") {
		t.Error("expected HasKey to return false for MISSING")
	}
}

func TestEqual_Identical(t *testing.T) {
	a := snapshot.New("secret/app", 1, baseSecrets(), snapshot.Meta{})
	b := snapshot.New("secret/app", 2, baseSecrets(), snapshot.Meta{})
	if !a.Equal(b) {
		t.Error("expected snapshots to be equal")
	}
}

func TestEqual_Different(t *testing.T) {
	a := snapshot.New("secret/app", 1, baseSecrets(), snapshot.Meta{})
	b := snapshot.New("secret/app", 2, map[string]string{"DB_HOST": "remotehost"}, snapshot.Meta{})
	if a.Equal(b) {
		t.Error("expected snapshots to be unequal")
	}
}
