package vault

import (
	"testing"
	"time"
)

func TestNewTokenSource_Valid(t *testing.T) {
	ts, err := NewTokenSource(TokenSourceDirect, "s.abc123", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.Token != "s.abc123" {
		t.Errorf("expected token s.abc123, got %s", ts.Token)
	}
	if ts.Type != TokenSourceDirect {
		t.Errorf("expected type direct, got %s", ts.Type)
	}
	if ts.ResolvedAt.IsZero() {
		t.Error("expected ResolvedAt to be set")
	}
}

func TestNewTokenSource_EmptyToken(t *testing.T) {
	_, err := NewTokenSource(TokenSourceEnv, "", 0)
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewTokenSource_EmptyType(t *testing.T) {
	_, err := NewTokenSource("", "s.abc123", 0)
	if err == nil {
		t.Fatal("expected error for empty type")
	}
}

func TestTokenSource_IsExpired_NoTTL(t *testing.T) {
	ts, _ := NewTokenSource(TokenSourceFile, "s.xyz", 0)
	if ts.IsExpired() {
		t.Error("token with zero TTL should never be expired")
	}
}

func TestTokenSource_IsExpired_Future(t *testing.T) {
	ts, _ := NewTokenSource(TokenSourceAppRole, "s.xyz", 10*time.Minute)
	if ts.IsExpired() {
		t.Error("token resolved now with 10m TTL should not be expired")
	}
}

func TestTokenSource_IsExpired_Past(t *testing.T) {
	ts := TokenSource{
		Type:       TokenSourceDirect,
		Token:      "s.old",
		ResolvedAt: time.Now().Add(-2 * time.Hour),
		TTL:        1 * time.Hour,
	}
	if !ts.IsExpired() {
		t.Error("expected token resolved 2h ago with 1h TTL to be expired")
	}
}

func TestTokenSource_String_NoExpiry(t *testing.T) {
	ts, _ := NewTokenSource(TokenSourceDirect, "s.tok", 0)
	s := ts.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestTokenSource_String_WithTTL(t *testing.T) {
	ts, _ := NewTokenSource(TokenSourceEnv, "s.tok", 30*time.Minute)
	s := ts.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestTokenSource_Validate_MissingResolvedAt(t *testing.T) {
	ts := TokenSource{
		Type:  TokenSourceDirect,
		Token: "s.abc",
	}
	if err := ts.Validate(); err == nil {
		t.Error("expected error for zero ResolvedAt")
	}
}
