// Package report provides report generation for vaultdiff.
//
// A Report aggregates diff results together with session metadata
// (user, timestamp, source/target paths) and can be rendered in
// multiple formats: plain text, JSON, and Markdown.
//
// Basic usage:
//
//	rep := report.New(session, results, "secret/dev", "secret/prod")
//	renderer := report.NewRenderer(report.FormatMarkdown, true)
//	renderer.Render(os.Stdout, rep)
package report
