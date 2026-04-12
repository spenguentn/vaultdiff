package vault

import (
	"testing"
	"time"
)

var baseLifecycle = SecretLifecycle{
	Mount:     "secret",
	Path:      "app/db",
	Stage:     LifecycleStageActive,
	ManagedBy: "platform-team",
}

func TestIsValidLifecycleStage_Known(t *testing.T) {
	for _, s := range []LifecycleStage{
		LifecycleStageActive, LifecycleStageDeprecated,
		LifecycleStageRetired, LifecycleStagePending,
	} {
		if !IsValidLifecycleStage(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidLifecycleStage_Unknown(t *testing.T) {
	if IsValidLifecycleStage("unknown") {
		t.Error("expected unknown stage to be invalid")
	}
}

func TestSecretLifecycle_FullPath(t *testing.T) {
	lc := baseLifecycle
	if got := lc.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretLifecycle_Validate_Valid(t *testing.T) {
	lc := baseLifecycle
	if err := lc.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretLifecycle_Validate_MissingMount(t *testing.T) {
	lc := baseLifecycle
	lc.Mount = ""
	if err := lc.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretLifecycle_Validate_InvalidStage(t *testing.T) {
	lc := baseLifecycle
	lc.Stage = "gone"
	if err := lc.Validate(); err == nil {
		t.Error("expected error for invalid stage")
	}
}

func TestSecretLifecycle_IsPastTransition_False(t *testing.T) {
	lc := baseLifecycle
	lc.TransitionAt = time.Now().Add(24 * time.Hour)
	if lc.IsPastTransition() {
		t.Error("expected IsPastTransition to be false")
	}
}

func TestSecretLifecycle_IsPastTransition_True(t *testing.T) {
	lc := baseLifecycle
	lc.TransitionAt = time.Now().Add(-1 * time.Hour)
	if !lc.IsPastTransition() {
		t.Error("expected IsPastTransition to be true")
	}
}

func TestLifecycleRegistry_SetAndGet(t *testing.T) {
	reg := NewSecretLifecycleRegistry()
	lc := baseLifecycle
	if err := reg.Set(&lc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("secret", "app/db")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.Stage != LifecycleStageActive {
		t.Errorf("unexpected stage: %s", got.Stage)
	}
}

func TestLifecycleRegistry_Set_SetsCreatedAt(t *testing.T) {
	reg := NewSecretLifecycleRegistry()
	lc := baseLifecycle
	_ = reg.Set(&lc)
	if lc.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestLifecycleRegistry_Set_Invalid(t *testing.T) {
	reg := NewSecretLifecycleRegistry()
	lc := baseLifecycle
	lc.ManagedBy = ""
	if err := reg.Set(&lc); err == nil {
		t.Error("expected validation error")
	}
}

func TestLifecycleRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretLifecycleRegistry()
	_, ok := reg.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestLifecycleRegistry_Remove(t *testing.T) {
	reg := NewSecretLifecycleRegistry()
	lc := baseLifecycle
	_ = reg.Set(&lc)
	if !reg.Remove("secret", "app/db") {
		t.Error("expected Remove to return true")
	}
	_, ok := reg.Get("secret", "app/db")
	if ok {
		t.Error("expected entry to be absent after removal")
	}
}

func TestLifecycleRegistry_All(t *testing.T) {
	reg := NewSecretLifecycleRegistry()
	a := baseLifecycle
	b := SecretLifecycle{Mount: "kv", Path: "svc/key", Stage: LifecycleStagePending, ManagedBy: "ops"}
	_ = reg.Set(&a)
	_ = reg.Set(&b)
	if len(reg.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(reg.All()))
	}
}
