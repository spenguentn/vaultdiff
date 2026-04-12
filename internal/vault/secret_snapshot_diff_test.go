package vault

import (
	"testing"
)

func TestSnapshotDiffEntry_FullPath(t *testing.T) {
	e := SnapshotDiffEntry{Mount: "secret", Path: "app/config", Key: "DB_PASS"}
	got := e.FullPath()
	want := "secret/app/config#DB_PASS"
	if got != want {
		t.Errorf("FullPath() = %q, want %q", got, want)
	}
}

func TestSnapshotDiffEntry_IsChanged_True(t *testing.T) {
	for _, ct := range []string{"added", "removed", "modified"} {
		e := SnapshotDiffEntry{ChangeType: ct}
		if !e.IsChanged() {
			t.Errorf("IsChanged() should be true for %q", ct)
		}
	}
}

func TestSnapshotDiffEntry_IsChanged_False(t *testing.T) {
	e := SnapshotDiffEntry{ChangeType: "unchanged"}
	if e.IsChanged() {
		t.Error("IsChanged() should be false for unchanged")
	}
}

func TestDiffSnapshots_Added(t *testing.T) {
	left := map[string]string{}
	right := map[string]string{"NEW_KEY": "value"}
	r := DiffSnapshots("left", "right", left, right)
	if len(r.Entries) != 1 || r.Entries[0].ChangeType != "added" {
		t.Errorf("expected 1 added entry, got %+v", r.Entries)
	}
}

func TestDiffSnapshots_Removed(t *testing.T) {
	left := map[string]string{"OLD_KEY": "value"}
	right := map[string]string{}
	r := DiffSnapshots("left", "right", left, right)
	if len(r.Entries) != 1 || r.Entries[0].ChangeType != "removed" {
		t.Errorf("expected 1 removed entry, got %+v", r.Entries)
	}
}

func TestDiffSnapshots_Modified(t *testing.T) {
	left := map[string]string{"KEY": "old"}
	right := map[string]string{"KEY": "new"}
	r := DiffSnapshots("left", "right", left, right)
	if len(r.Entries) != 1 || r.Entries[0].ChangeType != "modified" {
		t.Errorf("expected 1 modified entry, got %+v", r.Entries)
	}
}

func TestDiffSnapshots_Unchanged(t *testing.T) {
	left := map[string]string{"KEY": "same"}
	right := map[string]string{"KEY": "same"}
	r := DiffSnapshots("left", "right", left, right)
	if len(r.Entries) != 1 || r.Entries[0].ChangeType != "unchanged" {
		t.Errorf("expected 1 unchanged entry, got %+v", r.Entries)
	}
}

func TestSnapshotDiffResult_ChangedOnly(t *testing.T) {
	r := DiffSnapshots("l", "r",
		map[string]string{"A": "1", "B": "same"},
		map[string]string{"A": "2", "B": "same"},
	)
	changed := r.ChangedOnly()
	if len(changed) != 1 || changed[0].Key != "A" {
		t.Errorf("ChangedOnly() = %+v, want 1 modified entry for A", changed)
	}
}

func TestSnapshotDiffResult_Summary(t *testing.T) {
	r := DiffSnapshots("l", "r",
		map[string]string{"A": "old", "B": "same"},
		map[string]string{"A": "new", "B": "same", "C": "added"},
	)
	s := r.Summary()
	if s == "" {
		t.Error("Summary() should not be empty")
	}
}

func TestDiffSnapshots_Labels(t *testing.T) {
	r := DiffSnapshots("staging", "production", nil, nil)
	if r.LeftLabel != "staging" || r.RightLabel != "production" {
		t.Errorf("unexpected labels: %q / %q", r.LeftLabel, r.RightLabel)
	}
}
