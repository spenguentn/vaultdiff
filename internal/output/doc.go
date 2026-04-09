// Package output provides formatting and writing utilities for vaultdiff results.
//
// It supports multiple output formats including human-readable text and
// machine-readable JSON. Secret values can be masked based on a configurable
// list of sensitive key names before being written to the destination.
//
// Example usage:
//
//	w := output.NewWriter(output.FormatJSON, []string{"PASSWORD", "SECRET"}, os.Stdout)
//	if err := w.Write(results); err != nil {
//		log.Fatal(err)
//	}
package output
