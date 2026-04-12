package vault

import (
	"testing"
	"time"
)

var baseReminder = SecretReminder{
	Mount:     "secret",
	Path:      "app/db",
	Message:   "rotate this secret",
	RemindAt:  time.Now().Add(24 * time.Hour),
	Frequency: ReminderFrequencyOnce,
	CreatedBy: "alice",
}

func TestSecretReminder_FullPath(t *testing.T) {
	if got := baseReminder.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretReminder_Validate_Valid(t *testing.T) {
	if err := baseReminder.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestSecretReminder_Validate_MissingMount(t *testing.T) {
	r := baseReminder
	r.Mount = ""
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretReminder_Validate_MissingMessage(t *testing.T) {
	r := baseReminder
	r.Message = ""
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing message")
	}
}

func TestSecretReminder_Validate_ZeroRemindAt(t *testing.T) {
	r := baseReminder
	r.RemindAt = time.Time{}
	if err := r.Validate(); err == nil {
		t.Error("expected error for zero remind_at")
	}
}

func TestSecretReminder_IsDue_True(t *testing.T) {
	r := baseReminder
	r.RemindAt = time.Now().Add(-time.Minute)
	if !r.IsDue(time.Now()) {
		t.Error("expected reminder to be due")
	}
}

func TestSecretReminder_IsDue_False(t *testing.T) {
	r := baseReminder
	r.RemindAt = time.Now().Add(time.Hour)
	if r.IsDue(time.Now()) {
		t.Error("expected reminder to not be due")
	}
}

func TestIsValidFrequency_Known(t *testing.T) {
	for _, f := range []ReminderFrequency{ReminderFrequencyOnce, ReminderFrequencyWeekly, ReminderFrequencyMonthly} {
		if !IsValidFrequency(f) {
			t.Errorf("expected %s to be valid", f)
		}
	}
}

func TestIsValidFrequency_Unknown(t *testing.T) {
	if IsValidFrequency("daily") {
		t.Error("expected 'daily' to be invalid")
	}
}
