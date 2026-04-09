// Package filter provides utilities for narrowing diff results based on
// key prefixes, explicit key allowlists, and change-type predicates.
//
// Usage:
//
//	results := diff.Compare(left, right)
//	filtered := filter.Apply(results, filter.Options{
//		Prefix:      "db/",
//		OnlyChanged: true,
//	})
//
// Options can be combined; all active constraints must be satisfied for a
// result to be included in the output.
package filter
