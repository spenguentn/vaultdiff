package vault

import (
	"testing"
	"time"
)

func TestIsValidDispositionAction_Known(t *testing.T) {
	for _, a := range []DispositionAction{DispositionDelete, DispositionArchive, DispositionRotate, DispositionTransfer} {
		if !IsValidDispositionAction(a) {
			t.Errorf("expected %q to be valid", a)
		}
	}
}

func TestIsValidDispositionAction_Unknown(t *testing.T) {
	if IsValidDispositionAction("explode") {
		t.Error("expected unknown action to be invalid")
	}
}

func TestSecretDisposition_FullPath(t *testing.T) {
	d := &SecretDisposition{Mount: "secret", Path: "app/key"}
	if got := d.FullPath(); got != "secret/app/key" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretDisposition_IsDue_True(t *testing.T) {
	d := &SecretDisposition{ScheduledAt: time.Now().Add(-time.Hour)}
	if !d.IsDue(time.Now()) {
		t.Error("expected disposition to be due")
	}
}

func TestSecretDisposition_IsDue_False(t *testing.T) {
	d := &SecretDisposition{ScheduledAt: time.Now().Add(time.Hour)}
	if d.IsDue(time.Now()) {
		t.Error("expected disposition to not be due")
	}
}

func TestSecretDisposition_Validate_Valid(t *testing.T) {
	d := &SecretDisposition{
		Mount:       "secret",
		Path:        "app/db",
		Action:      DispositionDelete,
		ScheduledAt: time.Now().Add(24 * time.Hour),
		ApprovedBy:  "alice",
	}
	if err := d.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretDisposition_Validate_MissingMount(t *testing.T) {
	d := &SecretDisposition{Path: "app/db", Action: DispositionDelete, ScheduledAt: time.Now().Add(time.Hour), ApprovedBy: "alice"}
	if err := d.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretDisposition_Validate_InvalidAction(t *testing.T) {
	d := &SecretDisposition{Mount: "m", Path: "p", Action: "zap", ScheduledAt: time.Now().Add(time.Hour), ApprovedBy: "bob"}
	if err := d.Validate(); err == nil {
		t.Error("expected error for invalid action")
	}
}

func TestDispositionRegistry_SetAndGet(t *testing.T) {
	r := NewSecretDispositionRegistry()
	d := &SecretDisposition{Mount: "kv", Path: "svc/token", Action: DispositionArchive, ScheduledAt: time.Now().Add(time.Hour), ApprovedBy: "ops"}
	if err := r.Set(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get("kv", "svc/token")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if got.Action != DispositionArchive {
		t.Errorf("unexpected action: %s", got.Action)
	}
}

func TestDispositionRegistry_Due(t *testing.T) {
	r := NewSecretDispositionRegistry()
	past := &SecretDisposition{Mount: "kv", Path: "old", Action: DispositionDelete, ScheduledAt: time.Now().Add(-time.Minute), ApprovedBy: "admin"}
	future := &SecretDisposition{Mount: "kv", Path: "new", Action: DispositionRotate, ScheduledAt: time.Now().Add(time.Hour), ApprovedBy: "admin"}
	_ = r.Set(past)
	_ = r.Set(future)
	due := r.Due(time.Now())
	if len(due) != 1 {
		t.Errorf("expected 1 due disposition, got %d", len(due))
	}
}
