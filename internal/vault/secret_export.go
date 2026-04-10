package vault

import (
	"encoding/json"
	"fmt"
	"time"
)

// ExportFormat defines the serialization format for exported secrets.
type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatEnv  ExportFormat = "env"
)

// ExportOptions controls how secrets are exported.
type ExportOptions struct {
	Format     ExportFormat
	MaskValues bool
	IncludeMeta bool
}

// ExportRecord represents a single exported secret entry.
type ExportRecord struct {
	Path      string            `json:"path"`
	Version   int               `json:"version,omitempty"`
	ExportedAt time.Time        `json:"exported_at"`
	Data      map[string]string `json:"data"`
}

// Exporter serializes secret snapshots into portable formats.
type Exporter struct {
	opts ExportOptions
}

// NewExporter creates an Exporter with the given options.
func NewExporter(opts ExportOptions) *Exporter {
	if opts.Format == "" {
		opts.Format = ExportFormatJSON
	}
	return &Exporter{opts: opts}
}

// Export converts a map of secret data into an ExportRecord.
func (e *Exporter) Export(path string, version int, data map[string]string) ExportRecord {
	out := make(map[string]string, len(data))
	for k, v := range data {
		if e.opts.MaskValues {
			out[k] = "***"
		} else {
			out[k] = v
		}
	}
	return ExportRecord{
		Path:       path,
		Version:    version,
		ExportedAt: time.Now().UTC(),
		Data:       out,
	}
}

// MarshalJSON serializes an ExportRecord to JSON bytes.
func MarshalExportRecord(r ExportRecord) ([]byte, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("secret_export: marshal: %w", err)
	}
	return b, nil
}

// MarshalEnv serializes an ExportRecord as KEY=VALUE lines.
func MarshalEnv(r ExportRecord) []byte {
	var buf []byte
	for k, v := range r.Data {
		line := fmt.Sprintf("%s=%s\n", k, v)
		buf = append(buf, []byte(line)...)
	}
	return buf
}
