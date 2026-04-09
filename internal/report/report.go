// Package report provides structured report generation for vaultdiff results.
package report

import (
	"time"

	"github.com/vaultdiff/internal/audit"
	"github.com/vaultdiff/internal/diff"
)

// Format represents the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatMarkdown Format = "markdown"
)

// Report holds all data needed to render a diff report.
type Report struct {
	Session   *audit.Session
	Results   []diff.Result
	GeneratedAt time.Time
	SourcePath  string
	TargetPath  string
}

// New creates a new Report from the given session and diff results.
func New(session *audit.Session, results []diff.Result, src, target string) *Report {
	return &Report{
		Session:     session,
		Results:     results,
		GeneratedAt: time.Now().UTC(),
		SourcePath:  src,
		TargetPath:  target,
	}
}

// Summary returns a brief summary of the report's diff results.
func (r *Report) Summary() diff.Summary {
	return diff.Summary(r.Results)
}
