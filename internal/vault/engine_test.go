package vault

import (
	"testing"
)

func baseEngine() SecretsEngine {
	return SecretsEngine{
		Mount:       "secret",
		Type:        EngineKV2,
		Description: "default kv store",
	}
}

func TestSecretsEngine_Validate_Valid(t *testing.T) {
	e := baseEngine()
	if err := e.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSecretsEngine_Validate_EmptyMount(t *testing.T) {
	e := baseEngine()
	e.Mount = ""
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for empty mount, got nil")
	}
}

func TestSecretsEngine_Validate_EmptyType(t *testing.T) {
	e := baseEngine()
	e.Type = ""
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for empty type, got nil")
	}
}

func TestSecretsEngine_IsVersioned_True(t *testing.T) {
	e := baseEngine()
	e.Type = EngineKV2
	if !e.IsVersioned() {
		t.Fatal("expected IsVersioned to return true for kv-v2")
	}
}

func TestSecretsEngine_IsVersioned_False(t *testing.T) {
	e := baseEngine()
	e.Type = EngineKV1
	if e.IsVersioned() {
		t.Fatal("expected IsVersioned to return false for kv-v1")
	}
}

func TestSecretsEngine_String(t *testing.T) {
	e := baseEngine()
	got := e.String()
	want := "secret (kv-v2)"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestParseEngineType_KV1(t *testing.T) {
	if ParseEngineType("kv") != EngineKV1 {
		t.Fatal("expected EngineKV1 for 'kv'")
	}
	if ParseEngineType("kv-v1") != EngineKV1 {
		t.Fatal("expected EngineKV1 for 'kv-v1'")
	}
}

func TestParseEngineType_KV2(t *testing.T) {
	if ParseEngineType("kv-v2") != EngineKV2 {
		t.Fatal("expected EngineKV2 for 'kv-v2'")
	}
}

func TestParseEngineType_Unknown(t *testing.T) {
	if ParseEngineType("transit") != EngineUnknown {
		t.Fatal("expected EngineUnknown for unsupported type")
	}
}

func TestParseEngineType_Generic(t *testing.T) {
	if ParseEngineType("generic") != EngineGeneric {
		t.Fatal("expected EngineGeneric for 'generic'")
	}
}
