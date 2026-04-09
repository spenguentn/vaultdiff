package diff

import (
	"testing"
)

func TestCompare_Added(t *testing.T) {
	oldData := map[string]interface{}{}
	newData := map[string]interface{}{"token": "abc123"}

	result := Compare("secret/myapp", oldData, newData)

	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Change != Added {
		t.Errorf("expected Added, got %s", result.Entries[0].Change)
	}
	if !result.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestCompare_Removed(t *testing.T) {
	oldData := map[string]interface{}{"token": "abc123"}
	newData := map[string]interface{}{}

	result := Compare("secret/myapp", oldData, newData)

	if result.Entries[0].Change != Removed {
		t.Errorf("expected Removed, got %s", result.Entries[0].Change)
	}
	if result.Entries[0].OldValue != "abc123" {
		t.Errorf("unexpected OldValue: %s", result.Entries[0].OldValue)
	}
}

func TestCompare_Modified(t *testing.T) {
	oldData := map[string]interface{}{"password": "old"}
	newData := map[string]interface{}{"password": "new"}

	result := Compare("secret/myapp", oldData, newData)

	if result.Entries[0].Change != Modified {
		t.Errorf("expected Modified, got %s", result.Entries[0].Change)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	data := map[string]interface{}{"key": "value"}

	result := Compare("secret/myapp", data, data)

	if result.Entries[0].Change != Unchanged {
		t.Errorf("expected Unchanged, got %s", result.Entries[0].Change)
	}
	if result.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestCompare_SortedKeys(t *testing.T) {
	oldData := map[string]interface{}{"z": "1", "a": "2", "m": "3"}
	newData := map[string]interface{}{"z": "1", "a": "2", "m": "3"}

	result := Compare("secret/myapp", oldData, newData)

	if result.Entries[0].Key != "a" || result.Entries[1].Key != "m" || result.Entries[2].Key != "z" {
		t.Errorf("entries not sorted: %v", result.Entries)
	}
}

func TestMaskValue(t *testing.T) {
	if got := MaskValue("secret"); got != "s*****" {
		t.Errorf("unexpected mask: %s", got)
	}
	if got := MaskValue(""); got != "" {
		t.Errorf("unexpected mask for empty: %s", got)
	}
	if got := MaskValue("x"); got != "*" {
		t.Errorf("unexpected mask for single char: %s", got)
	}
}
