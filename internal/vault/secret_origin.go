package vault

import (
	"fmt"
	"strings"
	"time"
)

// OriginSource represents where a secret originated from.
type OriginSource string

const (
	OriginManual    OriginSource = "manual"
	OriginGenerated OriginSource = "generated"
	OriginImported  OriginSource = "imported"
	OriginReplicated OriginSource = "replicated"
)

// IsValidOriginSource returns true if the source is a known value.
func IsValidOriginSource(s OriginSource) bool {
	switch s {
	case OriginManual, OriginGenerated, OriginImported, OriginReplicated:
		return true
	}
	return false
}

// SecretOrigin records the provenance of a secret.
type SecretOrigin struct {
	Mount      string       `json:"mount"`
	Path       string       `json:"path"`
	Source     OriginSource `json:"source"`
	CreatedBy  string       `json:"created_by"`
	CreatedAt  time.Time    `json:"created_at"`
	ExternalRef string      `json:"external_ref,omitempty"`
}

// FullPath returns the combined mount and path.
func (o SecretOrigin) FullPath() string {
	return fmt.Sprintf("%s/%s", strings.Trim(o.Mount, "/"), strings.Trim(o.Path, "/"))
}

// Validate checks that the origin record has required fields.
func (o SecretOrigin) Validate() error {
	if o.Mount == "" {
		return fmt.Errorf("origin: mount is required")
	}
	if o.Path == "" {
		return fmt.Errorf("origin: path is required")
	}
	if o.CreatedBy == "" {
		return fmt.Errorf("origin: created_by is required")
	}
	if !IsValidOriginSource(o.Source) {
		return fmt.Errorf("origin: unknown source %q", o.Source)
	}
	return nil
}
