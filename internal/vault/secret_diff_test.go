package vault

import (
	"testing"
	"time"
)

func baseLeft() SecretVersion {
	return SecretVersion{
		Version:   1,
		CreatedAt: time.Now().Add(-time.Hour),
		Data: map[string]string{
			"host": "localhost",
			"port": "5432",
			"user": "admin",
		},
	}
}

func baseRight() SecretVersion {
	return SecretVersion{
		Version:   2,
		CreatedAt: time.Now(),
		Data: map[string]string{
			"host": "db.prod",
			"port": "5432",
			"pass": "s3cr3t",
		},
	}
}

func TestBuildSecretVersionDiff_HasChanges(t *testing.T) {
	d := BuildSecretVersionDiff("myapp/db", "secret", baseLeft(), baseRight())
	if !d.HasChanges() {
		t.Error("expected diff to have changes")
	}
}

func TestBuildSecretVersionDiff_Added(t *testing.T) {
	d := BuildSecretVersionDiff("myapp/db", "secret", baseLeft(), baseRight())
	if _, ok := d.Added["pass"]; !ok {
		t.Error("expected 'pass' to be in Added")
	}
}

func TestBuildSecretVersionDiff_Removed(t *testing.T) {
	d := BuildSecretVersionDiff("myapp/db", "secret", baseLeft(), baseRight())
	if _, ok := d.Removed["user"]; !ok {
		t.Error("expected 'user' to be in Removed")
	}
}

func TestBuildSecretVersionDiff_Modified(t *testing.T) {
	d := BuildSecretVersionDiff("myapp/db", "secret", baseLeft(), baseRight())
	if _, ok := d.Modified["host"]; !ok {
		t.Error("expected 'host' to be in Modified")
	}
}

func TestBuildSecretVersionDiff_Unchanged(t *testing.T) {
	d := BuildSecretVersionDiff("myapp/db", "secret", baseLeft(), baseRight())
	if _, ok := d.Unchanged["port"]; !ok {
		t.Error("expected 'port' to be in Unchanged")
	}
}

func TestBuildSecretVersionDiff_Versions(t *testing.T) {
	d := BuildSecretVersionDiff("myapp/db", "secret", baseLeft(), baseRight())
	if d.LeftVersion != 1 || d.RightVersion != 2 {
		t.Errorf("unexpected versions: left=%d right=%d", d.LeftVersion, d.RightVersion)
	}
}

func TestSecretVersionDiff_TotalKeys(t *testing.T) {
	d := BuildSecretVersionDiff("myapp/db", "secret", baseLeft(), baseRight())
	// added: pass, removed: user, modified: host, unchanged: port => 4
	if d.TotalKeys() != 4 {
		t.Errorf("expected 4 total keys, got %d", d.TotalKeys())
	}
}

func TestSecretVersionDiff_ChangedKeys(t *testing.T) {
	d := BuildSecretVersionDiff("myapp/db", "secret", baseLeft(), baseRight())
	// added: 1, removed: 1, modified: 1 => 3
	if d.ChangedKeys() != 3 {
		t.Errorf("expected 3 changed keys, got %d", d.ChangedKeys())
	}
}

func TestSecretVersionDiff_NoChanges(t *testing.T) {
	same := baseLeft()
	d := BuildSecretVersionDiff("myapp/db", "secret", same, same)
	if d.HasChanges() {
		t.Error("expected no changes when both versions are identical")
	}
}
