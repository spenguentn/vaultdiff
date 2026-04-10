package vault

import "time"

// SecretVersionDiff represents the difference between two versions of a secret.
type SecretVersionDiff struct {
	Path       string
	Mount      string
	LeftVersion  int
	RightVersion int
	LeftTime     time.Time
	RightTime    time.Time
	Added      map[string]string
	Removed    map[string]string
	Modified   map[string]string
	Unchanged  map[string]string
}

// HasChanges returns true if there are any added, removed, or modified keys.
func (d *SecretVersionDiff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Modified) > 0
}

// TotalKeys returns the total number of keys across all change categories.
func (d *SecretVersionDiff) TotalKeys() int {
	return len(d.Added) + len(d.Removed) + len(d.Modified) + len(d.Unchanged)
}

// ChangedKeys returns the count of keys that differ between versions.
func (d *SecretVersionDiff) ChangedKeys() int {
	return len(d.Added) + len(d.Removed) + len(d.Modified)
}

// BuildSecretVersionDiff constructs a SecretVersionDiff from two secret data maps
// and their associated metadata.
func BuildSecretVersionDiff(path, mount string, left, right SecretVersion) *SecretVersionDiff {
	added := make(map[string]string)
	removed := make(map[string]string)
	modified := make(map[string]string)
	unchanged := make(map[string]string)

	for k, rv := range right.Data {
		if lv, ok := left.Data[k]; !ok {
			added[k] = rv
		} else if lv != rv {
			modified[k] = rv
		} else {
			unchanged[k] = rv
		}
	}

	for k, lv := range left.Data {
		if _, ok := right.Data[k]; !ok {
			removed[k] = lv
		}
	}

	return &SecretVersionDiff{
		Path:         path,
		Mount:        mount,
		LeftVersion:  left.Version,
		RightVersion: right.Version,
		LeftTime:     left.CreatedAt,
		RightTime:    right.CreatedAt,
		Added:        added,
		Removed:      removed,
		Modified:     modified,
		Unchanged:    unchanged,
	}
}

// SecretVersion holds a parsed secret version and its metadata.
type SecretVersion struct {
	Version   int
	Data      map[string]string
	CreatedAt time.Time
}
