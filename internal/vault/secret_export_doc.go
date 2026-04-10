// Package vault provides secret export utilities for serializing Vault
// secret snapshots into portable formats such as JSON and env files.
//
// # Overview
//
// The export module consists of three main components:
//
//   - [Exporter]: converts raw secret data maps into [ExportRecord] values,
//     optionally masking sensitive values.
//
//   - [ExportWriter]: writes [ExportRecord] values to any [io.Writer] in the
//     configured [ExportFormat] (JSON or env).
//
//   - Marshal helpers: [MarshalExportRecord] and [MarshalEnv] provide
//     lower-level serialization without requiring an [ExportWriter].
//
// # Supported Formats
//
//   - ExportFormatJSON: pretty-printed JSON, suitable for audit logs.
//   - ExportFormatEnv:  KEY=VALUE lines, suitable for shell sourcing.
package vault
