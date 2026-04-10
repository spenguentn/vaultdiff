package vault

import (
	"errors"
	"fmt"
	"time"
)

// TokenSourceType identifies how a token is obtained.
type TokenSourceType string

const (
	TokenSourceDirect  TokenSourceType = "direct"
	TokenSourceEnv     TokenSourceType = "env"
	TokenSourceFile    TokenSourceType = "file"
	TokenSourceAppRole TokenSourceType = "approle"
)

// TokenSource describes where and how a Vault token was resolved.
type TokenSource struct {
	Type      TokenSourceType
	Token     string
	ResolvedAt time.Time
	TTL       n
// IsExpired reports whether the token has exceeded its TTL.
// If TTL is zero the token is considered non-expiring.
func (ts TokenSource) IsExpired() bool {
	if ts.TTL == 0 {
		return false
	}
	return time.Since(ts.ResolvedAt) > ts.TTL
}

// Validate checks that the TokenSource holds a usable token.
func (ts TokenSource) Validate() error {
	if ts.Token == "" {
		return errors.New("token source: token must not be empty")
	}
	if ts.Type == "" {
		return errors.New("token source: type must not be empty")
	}
	if ts.ResolvedAt.IsZero() {
		return errors.New("token source: resolved_at must be set")
	}
	return nil
}

// String returns a human-readable summary of the token source.
func (ts TokenSource) String() string {
	if ts.TTL == 0 {
		return fmt.Sprintf("TokenSource{type=%s, no_expiry}", ts.Type)
	}
	return fmt.Sprintf("TokenSource{type=%s, ttl=%s, expired=%v}", ts.Type, ts.TTL, ts.IsExpired())
}

// NewTokenSource constructs a TokenSource resolved at the current time.
func NewTokenSource(kind TokenSourceType, token string, ttl time.Duration) (TokenSource, error) {
	ts := TokenSource{
		Type:       kind,
		Token:      token,
		ResolvedAt: time.Now(),
		TTL:        ttl,
	}
	if err := ts.Validate(); err != nil {
		return TokenSource{}, err
	}
	return ts, nil
}
