// Package filter provides key-based filtering for secret diff results.
package filter

import (
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// Options controls which diff results are included in output.
type Options struct {
	// Prefix limits results to keys with this prefix.
	Prefix string
	// OnlyChanged excludes unchanged keys from results.
	OnlyChanged bool
	// Keys is an explicit allowlist of keys; empty means all keys.
	Keys []string
}

// Apply returns a filtered copy of results based on the given Options.
func Apply(results []diff.Result, opts Options) []diff.Result {
	allowlist := buildAllowlist(opts.Keys)

	filtered := make([]diff.Result, 0, len(results))
	for _, r := range results {
		if opts.OnlyChanged && r.Change == diff.Unchanged {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(r.Key, opts.Prefix) {
			continue
		}
		if len(allowlist) > 0 && !allowlist[r.Key] {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

// Count returns the number of results that would be returned by Apply.
// It is equivalent to len(Apply(results, opts)) but avoids allocating the
// result slice.
func Count(results []diff.Result, opts Options) int {
	allowlist := buildAllowlist(opts.Keys)

	count := 0
	for _, r := range results {
		if opts.OnlyChanged && r.Change == diff.Unchanged {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(r.Key, opts.Prefix) {
			continue
		}
		if len(allowlist) > 0 && !allowlist[r.Key] {
			continue
		}
		count++
	}
	return count
}

// buildAllowlist converts a slice of keys into a lookup map.
func buildAllowlist(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
