package vault

import "testing"

func sampleEngine(mount string, t EngineType) SecretsEngine {
	return SecretsEngine{Mount: mount, Type: t, Description: "test"}
}

func TestNewEngineRegistry_NotNil(t *testing.T) {
	r := NewEngineRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestEngineRegistry_Register_Valid(t *testing.T) {
	r := NewEngineRegistry()
	err := r.Register(sampleEngine("secret", EngineKV2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Count() != 1 {
		t.Fatalf("expected count 1, got %d", r.Count())
	}
}

func TestEngineRegistry_Register_Invalid(t *testing.T) {
	r := NewEngineRegistry()
	err := r.Register(SecretsEngine{Mount: "", Type: EngineKV2})
	if err == nil {
		t.Fatal("expected error for invalid engine, got nil")
	}
}

func TestEngineRegistry_Get_Found(t *testing.T) {
	r := NewEngineRegistry()
	_ = r.Register(sampleEngine("cubbyhole", EngineKV1))
	e, ok := r.Get("cubbyhole")
	if !ok {
		t.Fatal("expected engine to be found")
	}
	if e.Mount != "cubbyhole" {
		t.Fatalf("expected mount 'cubbyhole', got %q", e.Mount)
	}
}

func TestEngineRegistry_Get_NotFound(t *testing.T) {
	r := NewEngineRegistry()
	_, ok := r.Get("missing")
	if ok {
		t.Fatal("expected not found for unregistered mount")
	}
}

func TestEngineRegistry_Remove(t *testing.T) {
	r := NewEngineRegistry()
	_ = r.Register(sampleEngine("secret", EngineKV2))
	r.Remove("secret")
	if r.Count() != 0 {
		t.Fatalf("expected count 0 after remove, got %d", r.Count())
	}
}

func TestEngineRegistry_All_Count(t *testing.T) {
	r := NewEngineRegistry()
	_ = r.Register(sampleEngine("secret", EngineKV2))
	_ = r.Register(sampleEngine("cubbyhole", EngineKV1))
	if len(r.All()) != 2 {
		t.Fatalf("expected 2 engines, got %d", len(r.All()))
	}
}

func TestEngineRegistry_Register_Overwrite(t *testing.T) {
	r := NewEngineRegistry()
	_ = r.Register(sampleEngine("secret", EngineKV1))
	_ = r.Register(sampleEngine("secret", EngineKV2))
	e, _ := r.Get("secret")
	if e.Type != EngineKV2 {
		t.Fatalf("expected overwritten type kv-v2, got %q", e.Type)
	}
	if r.Count() != 1 {
		t.Fatalf("expected count 1 after overwrite, got %d", r.Count())
	}
}
