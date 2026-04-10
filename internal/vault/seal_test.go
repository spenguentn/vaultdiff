package vault

import (
	"strings"
	"testing"
)

func TestSealInfo_Status_Sealed(t *testing.T) {
	s := SealInfo{Initialized: true, Sealed: true}
	if s.Status() != SealStatusSealed {
		t.Fatalf("expected sealed, got %v", s.Status())
	}
}

func TestSealInfo_Status_Unsealed(t *testing.T) {
	s := SealInfo{Initialized: true, Sealed: false}
	if s.Status() != SealStatusUnsealed {
		t.Fatalf("expected unsealed, got %v", s.Status())
	}
}

func TestSealInfo_Status_Unknown(t *testing.T) {
	s := SealInfo{Initialized: false}
	if s.Status() != SealStatusUnknown {
		t.Fatalf("expected unknown, got %v", s.Status())
	}
}

func TestSealInfo_String_Sealed(t *testing.T) {
	s := SealInfo{Initialized: true, Sealed: true, Progress: 1, Threshold: 3}
	if !strings.Contains(s.String(), "sealed") {
		t.Fatalf("expected 'sealed' in string, got %q", s.String())
	}
}

func TestSealInfo_String_Unsealed(t *testing.T) {
	s := SealInfo{Initialized: true, Sealed: false, Version: "1.15.0"}
	if !strings.Contains(s.String(), "1.15.0") {
		t.Fatalf("expected version in string, got %q", s.String())
	}
}

func TestSealInfo_String_Unknown(t *testing.T) {
	s := SealInfo{Initialized: false}
	if !strings.Contains(s.String(), "unknown") {
		t.Fatalf("expected 'unknown' in string, got %q", s.String())
	}
}

func TestSealInfo_IsSealed_True(t *testing.T) {
	s := SealInfo{Sealed: true}
	if !s.IsSealed() {
		t.Fatal("expected IsSealed true")
	}
}

func TestSealInfo_IsSealed_False(t *testing.T) {
	s := SealInfo{Sealed: false}
	if s.IsSealed() {
		t.Fatal("expected IsSealed false")
	}
}

func TestParseSealInfo_Valid(t *testing.T) {
	raw := map[string]any{
		"sealed":      false,
		"initialized": true,
		"progress":    float64(0),
		"t":           float64(3),
		"n":           float64(5),
		"version":     "1.15.0",
	}
	info, err := ParseSealInfo(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Version != "1.15.0" {
		t.Errorf("expected version 1.15.0, got %q", info.Version)
	}
	if info.Threshold != 3 {
		t.Errorf("expected threshold 3, got %d", info.Threshold)
	}
	if info.Shares != 5 {
		t.Errorf("expected shares 5, got %d", info.Shares)
	}
}

func TestParseSealInfo_NilResponse(t *testing.T) {
	_, err := ParseSealInfo(nil)
	if err == nil {
		t.Fatal("expected error for nil response")
	}
}

func TestParseSealInfo_EmptyMap(t *testing.T) {
	info, err := ParseSealInfo(map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Sealed {
		t.Error("expected sealed=false for empty map")
	}
}
