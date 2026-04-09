package filter

import (
	"testing"
)

func TestOptions_IsZero_Empty(t *testing.T) {
	o := Options{}
	if !o.IsZero() {
		t.Error("expected IsZero() to return true for empty Options")
	}
}

func TestOptions_IsZero_WithOnlyChanged(t *testing.T) {
	o := Options{OnlyChanged: true}
	if o.IsZero() {
		t.Error("expected IsZero() to return false when OnlyChanged is set")
	}
}

func TestOptions_IsZero_WithPrefix(t *testing.T) {
	o := Options{Prefix: "app/"}
	if o.IsZero() {
		t.Error("expected IsZero() to return false when Prefix is set")
	}
}

func TestOptions_IsZero_WithMaxResults(t *testing.T) {
	o := Options{MaxResults: 10}
	if o.IsZero() {
		t.Error("expected IsZero() to return false when MaxResults is set")
	}
}

func TestOptions_HasKeyAllowlist_Empty(t *testing.T) {
	o := Options{}
	if o.HasKeyAllowlist() {
		t.Error("expected HasKeyAllowlist() to return false when allowlist is empty")
	}
}

func TestOptions_HasKeyAllowlist_Set(t *testing.T) {
	o := Options{KeyAllowlist: []string{"DB_HOST", "DB_PORT"}}
	if !o.HasKeyAllowlist() {
		t.Error("expected HasKeyAllowlist() to return true when keys are set")
	}
}

func TestOptions_HasExclusions_Empty(t *testing.T) {
	o := Options{}
	if o.HasExclusions() {
		t.Error("expected HasExclusions() to return false when no keys excluded")
	}
}

func TestOptions_HasExclusions_Set(t *testing.T) {
	o := Options{ExcludeKeys: []string{"SECRET_TOKEN"}}
	if !o.HasExclusions() {
		t.Error("expected HasExclusions() to return true when exclusions are set")
	}
}

func TestOptions_buildExcludeSet(t *testing.T) {
	o := Options{ExcludeKeys: []string{"KEY_A", "KEY_B", "KEY_C"}}
	set := o.buildExcludeSet()

	for _, k := range o.ExcludeKeys {
		if _, ok := set[k]; !ok {
			t.Errorf("expected key %q to be in exclude set", k)
		}
	}

	if len(set) != len(o.ExcludeKeys) {
		t.Errorf("expected set length %d, got %d", len(o.ExcludeKeys), len(set))
	}
}
