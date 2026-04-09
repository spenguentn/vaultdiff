// Package filter provides filtering capabilities for diff results.
package filter

// Options controls which diff results are included in output.
type Options struct {
	// OnlyChanged filters out unchanged entries when true.
	OnlyChanged bool

	// Prefix restricts results to keys matching the given prefix.
	Prefix string

	// KeyAllowlist restricts results to only the specified keys.
	// If empty, all keys are allowed (subject to other filters).
	KeyAllowlist []string

	// ExcludeKeys removes specific keys from the results.
	ExcludeKeys []string

	// MaxResults limits the number of results returned.
	// A value of 0 means no limit.
	MaxResults int
}

// IsZero returns true when no filtering options are set.
func (o Options) IsZero() bool {
	return !o.OnlyChanged &&
		o.Prefix == "" &&
		len(o.KeyAllowlist) == 0 &&
		len(o.ExcludeKeys) == 0 &&
		o.MaxResults == 0
}

// HasKeyAllowlist returns true when an allowlist of keys is configured.
func (o Options) HasKeyAllowlist() bool {
	return len(o.KeyAllowlist) > 0
}

// HasExclusions returns true when keys are explicitly excluded.
func (o Options) HasExclusions() bool {
	return len(o.ExcludeKeys) > 0
}

// buildExcludeSet converts ExcludeKeys into a set for O(1) lookup.
func (o Options) buildExcludeSet() map[string]struct{} {
	set := make(map[string]struct{}, len(o.ExcludeKeys))
	for _, k := range o.ExcludeKeys {
		set[k] = struct{}{}
	}
	return set
}
