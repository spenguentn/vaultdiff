package vault

import (
	"testing"
	"time"
)

func baseOrigin() SecretOrigin {
	return SecretOrigin{
		Mount:     "secret",
		Path:      "app/db",
		Source:    OriginManual,
		CreatedBy: "alice",
		CreatedAt: time.Now().UTC(),
	}
}

func TestIsValidOriginSource_Known(t *testing.T) {
	for _, s := range []OriginSource{OriginManual, OriginGenerated, OriginImported, OriginReplicated} {
		if !IsValidOriginSource(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidOriginSource_Unknown(t *testing.T) {
	if IsValidOriginSource("unknown") {
		t.Error("expected unknown source to be invalid")
	}
}

func TestSecretOrigin_FullPath(t *testing.T) {
	o := baseOrigin()
	if got := o.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretOrigin_Validate_Valid(t *testing.T) {
	if err := baseOrigin().Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretOrigin_Validate_MissingMount(t *testing.T) {
	o := baseOrigin()
	o.Mount = ""
	if err := o.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretOrigin_Validate_MissingCreatedBy(t *testing.T) {
	o := baseOrigin()
	o.CreatedBy = ""
	if err := o.Validate(); err == nil {
		t.Error("expected error for missing created_by")
	}
}

func TestSecretOrigin_Validate_InvalidSource(t *testing.T) {
	o := baseOrigin()
	o.Source = "unknown"
	if err := o.Validate(); err == nil {
		t.Error("expected error for invalid source")
	}
}

func TestNewSecretOriginRegistry_NotNil(t *testing.T) {
	if NewSecretOriginRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestOriginRegistry_Record_And_Get(t *testing.T) {
	r := NewSecretOriginRegistry()
	o := baseOrigin()
	if err := r.Record(o); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := r.Get(o.Mount, o.Path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.CreatedBy != o.CreatedBy {
		t.Errorf("expected %s, got %s", o.CreatedBy, got.CreatedBy)
	}
}

func TestOriginRegistry_Record_SetsCreatedAt(t *testing.T) {
	r := NewSecretOriginRegistry()
	o := baseOrigin()
	o.CreatedAt = time.Time{}
	_ = r.Record(o)
	got, _ := r.Get(o.Mount, o.Path)
	if got.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be stamped")
	}
}

func TestOriginRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretOriginRegistry()
	if _, err := r.Get("secret", "missing"); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestOriginRegistry_Remove(t *testing.T) {
	r := NewSecretOriginRegistry()
	o := baseOrigin()
	_ = r.Record(o)
	r.Remove(o.Mount, o.Path)
	if _, err := r.Get(o.Mount, o.Path); err == nil {
		t.Error("expected error after removal")
	}
}
