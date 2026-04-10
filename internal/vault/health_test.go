package vault

import (
	"testing"
	"time"
)

func TestHealthStatus_IsHealthy_True(t *testing.T) {
	h := HealthStatus{
		Initialized: true,
		Sealed:      false,
		Standby:     false,
	}
	if !h.IsHealthy() {
		t.Error("expected IsHealthy to return true")
	}
}

func TestHealthStatus_IsHealthy_Sealed(t *testing.T) {
	h := HealthStatus{
		Initialized: true,
		Sealed:      true,
		Standby:     false,
	}
	if h.IsHealthy() {
		t.Error("expected IsHealthy to return false when sealed")
	}
}

func TestHealthStatus_IsHealthy_NotInitialized(t *testing.T) {
	h := HealthStatus{
		Initialized: false,
		Sealed:      false,
		Standby:     false,
	}
	if h.IsHealthy() {
		t.Error("expected IsHealthy to return false when not initialized")
	}
}

func TestHealthStatus_IsHealthy_Standby(t *testing.T) {
	h := HealthStatus{
		Initialized: true,
		Sealed:      false,
		Standby:     true,
	}
	if h.IsHealthy() {
		t.Error("expected IsHealthy to return false when standby")
	}
}

func TestHealthStatus_String_Healthy(t *testing.T) {
	h := HealthStatus{
		Address:     "https://vault.example.com",
		Initialized: true,
		Sealed:      false,
		Standby:     false,
		Version:     "1.15.0",
	}
	got := h.String()
	want := "https://vault.example.com — healthy (v1.15.0)"
	if got != want {
		t.Errorf("String() = %q; want %q", got, want)
	}
}

func TestHealthStatus_String_Sealed(t *testing.T) {
	h := HealthStatus{
		Address:     "https://vault.example.com",
		Initialized: true,
		Sealed:      true,
	}
	got := h.String()
	want := "https://vault.example.com — sealed"
	if got != want {
		t.Errorf("String() = %q; want %q", got, want)
	}
}

func TestHealthStatus_CheckedAt_Set(t *testing.T) {
	before := time.Now().UTC()
	h := HealthStatus{CheckedAt: time.Now().UTC()}
	if h.CheckedAt.Before(before) {
		t.Error("expected CheckedAt to be set to a recent time")
	}
}

func TestCheckHealth_NilClient(t *testing.T) {
	_, err := CheckHealth(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil client, got nil")
	}
}
