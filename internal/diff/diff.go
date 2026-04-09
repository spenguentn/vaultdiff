package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ChangeType represents the type of change for a secret key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Entry represents a single key-level diff result.
type Entry struct {
	Key      string
	Change   ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two secret versions.
type Result struct {
	Path    string
	Entries []Entry
}

// HasChanges returns true if any entries differ.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Change != Unchanged {
			return true
		}
	}
	return false
}

// Compare diffs two secret data maps for a given path.
func Compare(path string, oldData, newData map[string]interface{}) *Result {
	result := &Result{Path: path}

	keys := unionKeys(oldData, newData)
	sort.Strings(keys)

	for _, k := range keys {
		oldVal, inOld := oldData[k]
		newVal, inNew := newData[k]

		entry := Entry{Key: k}

		switch {
		case inOld && !inNew:
			entry.Change = Removed
			entry.OldValue = fmt.Sprintf("%v", oldVal)
		case !inOld && inNew:
			entry.Change = Added
			entry.NewValue = fmt.Sprintf("%v", newVal)
		case fmt.Sprintf("%v", oldVal) != fmt.Sprintf("%v", newVal):
			entry.Change = Modified
			entry.OldValue = fmt.Sprintf("%v", oldVal)
			entry.NewValue = fmt.Sprintf("%v", newVal)
		default:
			entry.Change = Unchanged
			entry.OldValue = fmt.Sprintf("%v", oldVal)
			entry.NewValue = fmt.Sprintf("%v", newVal)
		}

		result.Entries = append(result.Entries, entry)
	}

	return result
}

func unionKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}

// MaskValue replaces all but the first character with asterisks.
func MaskValue(v string) string {
	if len(v) <= 1 {
		return strings.Repeat("*", len(v))
	}
	return string(v[0]) + strings.Repeat("*", len(v)-1)
}
