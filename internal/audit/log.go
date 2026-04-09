package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// Entry represents a single audit log record for a diff operation.
type Entry struct {
	Timestamp   time.Time        `json:"timestamp"`
	Environment string           `json:"environment"`
	Path        string           `json:"path"`
	FromVersion int              `json:"from_version"`
	ToVersion   int              `json:"to_version"`
	Changes     []diff.Change    `json:"changes"`
	User        string           `json:"user,omitempty"`
	Summary     diff.Summary     `json:"summary"`
}

// Logger writes structured audit entries to an io.Writer.
type Logger struct {
	w       io.Writer
	format  Format
}

// Format controls the output format of audit logs.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// NewLogger creates a new audit Logger writing to w.
func NewLogger(w io.Writer, format Format) *Logger {
	return &Logger{w: w, format: format}
}

// Write records an audit entry to the underlying writer.
func (l *Logger) Write(entry Entry) error {
	switch l.format {
	case FormatJSON:
		return l.writeJSON(entry)
	case FormatText:
		return l.writeText(entry)
	default:
		return fmt.Errorf("unsupported audit format: %s", l.format)
	}
}

func (l *Logger) writeJSON(entry Entry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: failed to marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", data)
	return err
}

func (l *Logger) writeText(entry Entry) error {
	_, err := fmt.Fprintf(l.w,
		"[%s] env=%s path=%s versions=%d->%d added=%d removed=%d modified=%d user=%s\n",
		entry.Timestamp.UTC().Format(time.RFC3339),
		entry.Environment,
		entry.Path,
		entry.FromVersion,
		entry.ToVersion,
		entry.Summary.Added,
		entry.Summary.Removed,
		entry.Summary.Modified,
		entry.User,
	)
	return err
}
