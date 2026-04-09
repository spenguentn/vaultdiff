package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/vaultdiff/internal/diff"
)

// Renderer writes a Report to an io.Writer in a given Format.
type Renderer struct {
	format      Format
	maskSecrets bool
}

// NewRenderer creates a Renderer for the specified format.
func NewRenderer(format Format, maskSecrets bool) *Renderer {
	return &Renderer{format: format, maskSecrets: maskSecrets}
}

// Render writes the report to w.
func (r *Renderer) Render(w io.Writer, rep *Report) error {
	switch r.format {
	case FormatJSON:
		return r.renderJSON(w, rep)
	case FormatMarkdown:
		return r.renderMarkdown(w, rep)
	case FormatText:
		return r.renderText(w, rep)
	default:
		return fmt.Errorf("unsupported report format: %s", r.format)
	}
}

type jsonReport struct {
	GeneratedAt time.Time    `json:"generated_at"`
	Source      string       `json:"source"`
	Target      string       `json:"target"`
	Results     []diff.Result `json:"results"`
}

func (r *Renderer) renderJSON(w io.Writer, rep *Report) error {
	results := applyMask(rep.Results, r.maskSecrets)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonReport{
		GeneratedAt: rep.GeneratedAt,
		Source:      rep.SourcePath,
		Target:      rep.TargetPath,
		Results:     results,
	})
}

func (r *Renderer) renderText(w io.Writer, rep *Report) error {
	fmt.Fprintf(w, "VaultDiff Report — %s\n", rep.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Source : %s\n", rep.SourcePath)
	fmt.Fprintf(w, "Target : %s\n", rep.TargetPath)
	fmt.Fprintln(w, strings.Repeat("-", 60))
	for _, res := range applyMask(rep.Results, r.maskSecrets) {
		fmt.Fprintf(w, "[%s] %s\n", res.ChangeType, res.Key)
	}
	s := rep.Summary()
	fmt.Fprintf(w, "\nSummary: +%d added  ~%d modified  -%d removed  =%d unchanged\n",
		s.Added, s.Modified, s.Removed, s.Unchanged)
	return nil
}

func (r *Renderer) renderMarkdown(w io.Writer, rep *Report) error {
	fmt.Fprintf(w, "# VaultDiff Report\n\n")
	fmt.Fprintf(w, "**Generated:** %s  \n", rep.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "**Source:** `%s`  \n", rep.SourcePath)
	fmt.Fprintf(w, "**Target:** `%s`  \n\n", rep.TargetPath)
	fmt.Fprintln(w, "| Change | Key | Old Value | New Value |")
	fmt.Fprintln(w, "|--------|-----|-----------|-----------|")
	for _, res := range applyMask(rep.Results, r.maskSecrets) {
		fmt.Fprintf(w, "| %s | `%s` | %s | %s |\n", res.ChangeType, res.Key, res.OldValue, res.NewValue)
	}
	return nil
}

func applyMask(results []diff.Result, mask bool) []diff.Result {
	if !mask {
		return results
	}
	out := make([]diff.Result, len(results))
	for i, r := range results {
		out[i] = r
		out[i].OldValue = diff.MaskValue(r.OldValue)
		out[i].NewValue = diff.MaskValue(r.NewValue)
	}
	return out
}
