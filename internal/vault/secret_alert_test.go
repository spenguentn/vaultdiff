package vault

import (
	"testing"
	"time"
)

func baseAlert() SecretAlert {
	return SecretAlert{
		Mount:    "secret",
		Path:     "app/db",
		Message:  "secret is about to expire",
		Severity: AlertSeverityHigh,
		Actor:    "ci-bot",
	}
}

func TestSecretAlert_FullPath(t *testing.T) {
	a := baseAlert()
	if got := a.FullPath(); got != "secret/app/db" {
		t.Fatalf("expected secret/app/db, got %s", got)
	}
}

func TestSecretAlert_Validate_Valid(t *testing.T) {
	if err := baseAlert().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretAlert_Validate_MissingMount(t *testing.T) {
	a := baseAlert()
	a.Mount = ""
	if err := a.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestSecretAlert_Validate_MissingMessage(t *testing.T) {
	a := baseAlert()
	a.Message = ""
	if err := a.Validate(); err == nil {
		t.Fatal("expected error for missing message")
	}
}

func TestSecretAlert_Validate_InvalidSeverity(t *testing.T) {
	a := baseAlert()
	a.Severity = "unknown"
	if err := a.Validate(); err == nil {
		t.Fatal("expected error for invalid severity")
	}
}

func TestIsValidSeverity_Known(t *testing.T) {
	for _, s := range []AlertSeverity{AlertSeverityLow, AlertSeverityMedium, AlertSeverityHigh, AlertSeverityCritical} {
		if !IsValidSeverity(s) {
			t.Fatalf("expected %s to be valid", s)
		}
	}
}

func TestIsValidSeverity_Unknown(t *testing.T) {
	if IsValidSeverity("extreme") {
		t.Fatal("expected extreme to be invalid")
	}
}

func TestAlertRegistry_Record_And_Get(t *testing.T) {
	r := NewSecretAlertRegistry()
	a := baseAlert()
	a.Triggered = time.Now().UTC()
	if err := r.Record(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, ok := r.Get("secret", "app/db")
	if !ok || len(list) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(list))
	}
}

func TestAlertRegistry_Record_SetsTriggered(t *testing.T) {
	r := NewSecretAlertRegistry()
	a := baseAlert()
	_ = r.Record(a)
	list, _ := r.Get("secret", "app/db")
	if list[0].Triggered.IsZero() {
		t.Fatal("expected Triggered to be set automatically")
	}
}

func TestAlertRegistry_Clear(t *testing.T) {
	r := NewSecretAlertRegistry()
	a := baseAlert()
	a.Triggered = time.Now().UTC()
	_ = r.Record(a)
	r.Clear("secret", "app/db")
	_, ok := r.Get("secret", "app/db")
	if ok {
		t.Fatal("expected alerts to be cleared")
	}
}

func TestAlertRegistry_All(t *testing.T) {
	r := NewSecretAlertRegistry()
	a1 := baseAlert()
	a1.Triggered = time.Now().UTC()
	a2 := baseAlert()
	a2.Path = "app/cache"
	a2.Triggered = time.Now().UTC()
	_ = r.Record(a1)
	_ = r.Record(a2)
	if got := len(r.All()); got != 2 {
		t.Fatalf("expected 2 alerts, got %d", got)
	}
}
