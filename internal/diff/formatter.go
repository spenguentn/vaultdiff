package diff

import (
	"fmt"
	"io"
	"strings"
)

// Format controls how diff output is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// TextFormatter writes a human-readable diff to the given writer.
type TextFormatter struct {
	MaskSecrets bool
	Writer      io.Writer
}

// Write renders the diff result as colored text output.
func (f *TextFormatter) Write(result *Result) error {
	if !result.HasChanges() {
		fmt.Fprintf(f.Writer, "[=] %s — no changes\n", result.Path)
		return nil
	}

	fmt.Fprintf(f.Writer, "[~] %s\n", result.Path)
	fmt.Fprintf(f.Writer, "%s\n", strings.Repeat("-", 40))

	for _, e := range result.Entries {
		switch e.Change {
		case Added:
			newVal := e.NewValue
			if f.MaskSecrets {
				newVal = MaskValue(newVal)
			}
			fmt.Fprintf(f.Writer, "  + %-20s %s\n", e.Key, newVal)
		case Removed:
			oldVal := e.OldValue
			if f.MaskSecrets {
				oldVal = MaskValue(oldVal)
			}
			fmt.Fprintf(f.Writer, "  - %-20s %s\n", e.Key, oldVal)
		case Modified:
			oldVal, newVal := e.OldValue, e.NewValue
			if f.MaskSecrets {
				oldVal = MaskValue(oldVal)
				newVal = MaskValue(newVal)
			}
			fmt.Fprintf(f.Writer, "  ~ %-20s %s → %s\n", e.Key, oldVal, newVal)
		}
	}

	fmt.Fprintln(f.Writer)
	return nil
}

// Summary returns a one-line summary of changes.
func Summary(result *Result) string {
	var added, removed, modified int
	for _, e := range result.Entries {
		switch e.Change {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return fmt.Sprintf("%s: +%d -%d ~%d", result.Path, added, removed, modified)
}
