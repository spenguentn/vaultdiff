package vault

import (
	"fmt"
	"time"
)

// SnapshotDiffEntry represents a single key-level change between two snapshots.
type SnapshotDiffEntry struct {
	Mount     string
	Path      string
	Key       string
	ChangeType string // "added", "removed", "modified", "unchanged"
	LeftValue  string
	RightValue string
	RecordedAt time.Time
}

// FullPath returns the canonical path for the entry.
func (e SnapshotDiffEntry) FullPath() string {
	return fmt.Sprintf("%s/%s#%s", e.Mount, e.Path, e.Key)
}

// IsChanged returns true if the entry reflects an actual change.
func (e SnapshotDiffEntry) IsChanged() bool {
	return e.ChangeType != "unchanged"
}

// SnapshotDiffResult holds the full result of diffing two snapshots.
type SnapshotDiffResult struct {
	LeftLabel  string
	RightLabel string
	Entries    []SnapshotDiffEntry
	GeneratedAt time.Time
}

// ChangedOnly returns only entries that represent a change.
func (r *SnapshotDiffResult) ChangedOnly() []SnapshotDiffEntry {
	out := make([]SnapshotDiffEntry, 0, len(r.Entries))
	for _, e := range r.Entries {
		if e.IsChanged() {
			out = append(out, e)
		}
	}
	return out
}

// Summary returns a short human-readable summary of the diff.
func (r *SnapshotDiffResult) Summary() string {
	var added, removed, modified, unchanged int
	for _, e := range r.Entries {
		switch e.ChangeType {
		case "added":
			added++
		case "removed":
			removed++
		case "modified":
			modified++
		default:
			unchanged++
		}
	}
	return fmt.Sprintf("added=%d removed=%d modified=%d unchanged=%d", added, removed, modified, unchanged)
}

// DiffSnapshots compares two maps of secret data and returns a SnapshotDiffResult.
func DiffSnapshots(leftLabel, rightLabel string, left, right map[string]string) *SnapshotDiffResult {
	result := &SnapshotDiffResult{
		LeftLabel:   leftLabel,
		RightLabel:  rightLabel,
		GeneratedAt: time.Now().UTC(),
	}

	seen := make(map[string]bool)
	for k, lv := range left {
		seen[k] = true
		rv, ok := right[k]
		ct := "unchanged"
		if !ok {
			ct = "removed"
		} else if lv != rv {
			ct = "modified"
		}
		result.Entries = append(result.Entries, SnapshotDiffEntry{
			Key:        k,
			ChangeType: ct,
			LeftValue:  lv,
			RightValue: rv,
			RecordedAt: result.GeneratedAt,
		})
	}
	for k, rv := range right {
		if !seen[k] {
			result.Entries = append(result.Entries, SnapshotDiffEntry{
				Key:        k,
				ChangeType: "added",
				RightValue: rv,
				RecordedAt: result.GeneratedAt,
			})
		}
	}
	return result
}
