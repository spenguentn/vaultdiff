package vault

import (
	"regexp"
	"strings"
)

// RedactMode controls how sensitive values are redacted.
type RedactMode string

const (
	RedactMask    RedactMode = "mask"    // Replace with fixed mask string
	RedactPartial RedactMode = "partial" // Show first/last N chars
	RedactRemove  RedactMode = "remove"  // Remove the key entirely
)

// RedactConfig defines rules for redacting secret values.
type RedactConfig struct {
	Mode        RedactMode
	MaskString  string
	Patterns    []*regexp.Regexp
	KeyPrefixes []string
}

// DefaultRedactConfig returns a RedactConfig with sensible defaults.
func DefaultRedactConfig() RedactConfig {
	return RedactConfig{
		Mode:       RedactMask,
		MaskString: "[REDACTED]",
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)password`),
			regexp.MustCompile(`(?i)secret`),
			regexp.MustCompile(`(?i)token`),
			regexp.MustCompile(`(?i)api[_-]?key`),
			regexp.MustCompile(`(?i)private[_-]?key`),
		},
	}
}

// ShouldRedact reports whether the given key matches any redaction rule.
func (c RedactConfig) ShouldRedact(key string) bool {
	for _, prefix := range c.KeyPrefixes {
		if strings.HasPrefix(strings.ToLower(key), strings.ToLower(prefix)) {
			return true
		}
	}
	for _, re := range c.Patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// Apply redacts values in the provided map according to the config.
// Returns a new map; the original is not modified.
func (c RedactConfig) Apply(data map[string]string) map[string]string {
	out := make(map[string]string, len(data))
	for k, v := range data {
		if !c.ShouldRedact(k) {
			out[k] = v
			continue
		}
		switch c.Mode {
		case RedactPartial:
			out[k] = partialMask(v)
		case RedactRemove:
			// omit key entirely
		default:
			mask := c.MaskString
			if mask == "" {
				mask = "[REDACTED]"
			}
			out[k] = mask
		}
	}
	return out
}

// partialMask shows the first and last character of a value if long enough.
func partialMask(v string) string {
	if len(v) <= 2 {
		return strings.Repeat("*", len(v))
	}
	return string(v[0]) + strings.Repeat("*", len(v)-2) + string(v[len(v)-1])
}
