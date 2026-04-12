package vault

import (
	"testing"
)

func TestComputeChecksum_Valid(t *testing.T) {
	c, err := ComputeChecksum("secret", "myapp/db", 1, map[string]interface{}{
		"password": "s3cr3t",
		"user":     "admin",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
	if c.Mount != "secret" || c.Path != "myapp/db" {
		t.Errorf("unexpected mount/path: %s/%s", c.Mount, c.Path)
	}
}

func TestComputeChecksum_Deterministic(t *testing.T) {
	data := map[string]interface{}{"b": "2", "a": "1"}
	c1, _ := ComputeChecksum("kv", "app", 1, data)
	c2, _ := ComputeChecksum("kv", "app", 1, data)
	if c1.Checksum != c2.Checksum {
		t.Error("checksums should be deterministic")
	}
}

func TestComputeChecksum_MissingMount(t *testing.T) {
	_, err := ComputeChecksum("", "app", 1, nil)
	if err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretChecksum_Matches_True(t *testing.T) {
	c1, _ := ComputeChecksum("kv", "app", 1, map[string]interface{}{"key": "val"})
	c2, _ := ComputeChecksum("kv", "app", 1, map[string]interface{}{"key": "val"})
	if !c1.Matches(c2) {
		t.Error("expected checksums to match")
	}
}

func TestSecretChecksum_Matches_False(t *testing.T) {
	c1, _ := ComputeChecksum("kv", "app", 1, map[string]interface{}{"key": "val"})
	c2, _ := ComputeChecksum("kv", "app", 2, map[string]interface{}{"key": "changed"})
	if c1.Matches(c2) {
		t.Error("expected checksums not to match")
	}
}

func TestNewSecretChecksumRegistry_NotNil(t *testing.T) {
	r := NewSecretChecksumRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestChecksumRegistry_StoreAndGet(t *testing.T) {
	r := NewSecretChecksumRegistry()
	c, _ := ComputeChecksum("kv", "app", 1, map[string]interface{}{"x": "y"})
	if err := r.Store(c); err != nil {
		t.Fatalf("unexpected store error: %v", err)
	}
	got, err := r.Get("kv", "app")
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}
	if got.Checksum != c.Checksum {
		t.Errorf("checksum mismatch: got %s want %s", got.Checksum, c.Checksum)
	}
}

func TestChecksumRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretChecksumRegistry()
	_, err := r.Get("kv", "missing")
	if err == nil {
		t.Error("expected error for missing entry")
	}
}

func TestChecksumRegistry_Remove(t *testing.T) {
	r := NewSecretChecksumRegistry()
	c, _ := ComputeChecksum("kv", "app", 1, map[string]interface{}{"k": "v"})
	_ = r.Store(c)
	r.Remove("kv", "app")
	if r.Len() != 0 {
		t.Error("expected registry to be empty after remove")
	}
}

func TestChecksumRegistry_All(t *testing.T) {
	r := NewSecretChecksumRegistry()
	c1, _ := ComputeChecksum("kv", "a", 1, map[string]interface{}{})
	c2, _ := ComputeChecksum("kv", "b", 1, map[string]interface{}{})
	_ = r.Store(c1)
	_ = r.Store(c2)
	if len(r.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.All()))
	}
}
