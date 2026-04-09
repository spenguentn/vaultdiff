package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/diff"
)

// Format represents the output format for diff results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Writer writes diff results to an output destination.
type Writer struct {
	format    Format
	maskKeys  []string
	dest      io.Writer
	formatter *diff.TextFormatter
}

// NewWriter creates a new Writer with the given format and destination.
func NewWriter(format Format, maskKeys []string, dest io.Writer) *Writer {
	if dest == nil {
		dest = os.Stdout
	}
	return &Writer{
		format:    format,
		maskKeys:  maskKeys,
		dest:      dest,
		formatter: diff.NewTextFormatter(maskKeys),
	}
}

// Write outputs the diff results according to the configured format.
func (w *Writer) Write(results []diff.Result) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(results)
	case FormatText:
		return w.writeText(results)
	default:
		return fmt.Errorf("unsupported output format: %s", w.format)
	}
}

func (w *Writer) writeText(results []diff.Result) error {
	output := w.formatter.Format(results)
	_, err := fmt.Fprint(w.dest, output)
	return err
}

func (w *Writer) writeJSON(results []diff.Result) error {
	type jsonResult struct {
		Key      string `json:"key"`
		ChangeType string `json:"change_type"`
		OldValue string `json:"old_value,omitempty"`
		NewValue string `json:"new_value,omitempty"`
	}

	out := make([]jsonResult, 0, len(results))
	for _, r := range results {
		out = append(out, jsonResult{
			Key:        r.Key,
			ChangeType: string(r.ChangeType),
			OldValue:   diff.MaskValue(r.Key, r.OldValue, w.maskKeys),
			NewValue:   diff.MaskValue(r.Key, r.NewValue, w.maskKeys),
		})
	}

	enc := json.NewEncoder(w.dest)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
