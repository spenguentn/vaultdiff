package vault

import (
	"testing"
)

func TestMountInfo_Validate_Valid(t *testing.T) {
	m := MountInfo{Path: "secret", Type: MountTypeKV2}
	if err := m.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestMountInfo_Validate_EmptyPath(t *testing.T) {
	m := MountInfo{Path: "", Type: MountTypeKV2}
	if err := m.Validate(); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestMountInfo_Validate_EmptyType(t *testing.T) {
	m := MountInfo{Path: "secret", Type: ""}
	if err := m.Validate(); err == nil {
		t.Fatal("expected error for empty type")
	}
}

func TestMountInfo_IsKV_V1(t *testing.T) {
	m := MountInfo{Path: "secret", Type: MountTypeKV1}
	if !m.IsKV() {
		t.Fatal("expected IsKV to be true for kv-v1")
	}
}

func TestMountInfo_IsKV_V2(t *testing.T) {
	m := MountInfo{Path: "secret", Type: MountTypeKV2}
	if !m.IsKV() {
		t.Fatal("expected IsKV to be true for kv-v2")
	}
}

func TestMountInfo_IsKV_NonKV(t *testing.T) {
	m := MountInfo{Path: "pki", Type: MountTypePKI}
	if m.IsKV() {
		t.Fatal("expected IsKV to be false for pki")
	}
}

func TestParseMountType_KVAlias(t *testing.T) {
	if got := ParseMountType("kv"); got != MountTypeKV1 {
		t.Fatalf("expected %q, got %q", MountTypeKV1, got)
	}
}

func TestParseMountType_KV2(t *testing.T) {
	if got := ParseMountType("kv-v2"); got != MountTypeKV2 {
		t.Fatalf("expected %q, got %q", MountTypeKV2, got)
	}
}

func TestParseMountType_CaseInsensitive(t *testing.T) {
	if got := ParseMountType("PKI"); got != MountTypePKI {
		t.Fatalf("expected %q, got %q", MountTypePKI, got)
	}
}

func TestParseMountType_Unknown(t *testing.T) {
	raw := "custom-engine"
	if got := ParseMountType(raw); string(got) != raw {
		t.Fatalf("expected %q passthrough, got %q", raw, got)
	}
}
