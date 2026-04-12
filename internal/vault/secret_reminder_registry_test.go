package vault

import (
	"testing"
	"time"
)

func TestNewSecretReminderRegistry_NotNil(t *testing.T) {
	if NewSecretReminderRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestReminderRegistry_Set_And_Get(t *testing.T) {
	reg := NewSecretReminderRegistry()
	if err := reg.Set(baseReminder); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := reg.Get(baseReminder.Mount, baseReminder.Path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Message != baseReminder.Message {
		t.Errorf("message mismatch: got %s", got.Message)
	}
}

func TestReminderRegistry_Set_SetsCreatedAt(t *testing.T) {
	reg := NewSecretReminderRegistry()
	r := baseReminder
	r.CreatedAt = time.Time{}
	_ = reg.Set(r)
	got, _ := reg.Get(r.Mount, r.Path)
	if got.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be stamped")
	}
}

func TestReminderRegistry_Set_Invalid(t *testing.T) {
	reg := NewSecretReminderRegistry()
	r := baseReminder
	r.Mount = ""
	if err := reg.Set(r); err == nil {
		t.Error("expected validation error")
	}
}

func TestReminderRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretReminderRegistry()
	if _, err := reg.Get("secret", "missing"); err == nil {
		t.Error("expected not-found error")
	}
}

func TestReminderRegistry_Remove(t *testing.T) {
	reg := NewSecretReminderRegistry()
	_ = reg.Set(baseReminder)
	if err := reg.Remove(baseReminder.Mount, baseReminder.Path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := reg.Get(baseReminder.Mount, baseReminder.Path); err == nil {
		t.Error("expected not-found after remove")
	}
}

func TestReminderRegistry_Remove_NotFound(t *testing.T) {
	reg := NewSecretReminderRegistry()
	if err := reg.Remove("secret", "ghost"); err == nil {
		t.Error("expected error removing non-existent reminder")
	}
}

func TestReminderRegistry_Due(t *testing.T) {
	reg := NewSecretReminderRegistry()
	past := baseReminder
	past.Path = "app/past"
	past.RemindAt = time.Now().Add(-time.Minute)

	future := baseReminder
	future.Path = "app/future"
	future.RemindAt = time.Now().Add(time.Hour)

	_ = reg.Set(past)
	_ = reg.Set(future)

	due := reg.Due(time.Now())
	if len(due) != 1 {
		t.Errorf("expected 1 due reminder, got %d", len(due))
	}
	if due[0].Path != "app/past" {
		t.Errorf("unexpected due path: %s", due[0].Path)
	}
}
