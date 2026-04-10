package vault

import (
	"fmt"
	"io"
)

// ExportWriter writes ExportRecords to an io.Writer in the configured format.
type ExportWriter struct {
	w    io.Writer
	opts ExportOptions
}

// NewExportWriter creates an ExportWriter targeting w.
func NewExportWriter(w io.Writer, opts ExportOptions) *ExportWriter {
	if opts.Format == "" {
		opts.Format = ExportFormatJSON
	}
	return &ExportWriter{w: w, opts: opts}
}

// Write serializes rec to the underlying writer.
func (ew *ExportWriter) Write(rec ExportRecord) error {
	switch ew.opts.Format {
	case ExportFormatJSON:
		b, err := MarshalExportRecord(rec)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(ew.w, string(b))
		return err
	case ExportFormatEnv:
		_, err := ew.w.Write(MarshalEnv(rec))
		return err
	default:
		return fmt.Errorf("secret_export_writer: unsupported format %q", ew.opts.Format)
	}
}

// WriteAll writes multiple records sequentially.
func (ew *ExportWriter) WriteAll(records []ExportRecord) error {
	for _, r := range records {
		if err := ew.Write(r); err != nil {
			return fmt.Errorf("secret_export_writer: write all: %w", err)
		}
	}
	return nil
}
