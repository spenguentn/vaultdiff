package vault

import (
	"testing"
	"time"
)

var baseAccessEntry = SecretAccessEntry{
	Mount:     "secret",
	Path:      "app/db",
	EventType: AccessEventRead,
	Actor:     "ci-bot",
	Timestamp: time.Now().UTC(),
}

func TestSecretAccessEntry_FullPath(t *testing.T) {
	e := baseAccessEntry
	if got := e.FullPath(); got != "secret/app/db" {
		t.Errorf("expected secret/app/db, got %s", got)
	}
}

func TestSecretAccessEntry_Validate_Valid(t *testing.T) {
	e := baseAccessEntry
	if err := e.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretAccessEntry_Validate_MissingMount(t *testing.T) {
	e := baseAccessEntry
	e.Mount = ""
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretAccessEntry_Validate_MissingActor(t *testing.T) {
	e := baseAccessEntry
	e.Actor = ""
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing actor")
	}
}

func TestIsValidEventType_Known(t *testing.T) {
	for _, et := range []AccessEventType{AccessEventRead, AccessEventWrite, AccessEventDelete, AccessEventList} {
		if !IsValidEventType(et) {
			t.Errorf("expected %s to be valid", et)
		}
	}
}

func TestIsValidEventType_Unknown(t *testing.T) {
	if IsValidEventType("unknown") {
		t.Error("expected unknown to be invalid")
	}
}
