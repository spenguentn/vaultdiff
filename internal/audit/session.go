package audit

import (
	"os/user"
	"time"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// Session captures context for a single diff invocation and
// produces an Entry ready for the Logger.
type Session struct {
	Environment string
	Path        string
	FromVersion int
	ToVersion   int
	User        string
}

// NewSession creates a Session, auto-detecting the OS user if user is empty.
func NewSession(env, path string, from, to int, userName string) (*Session, error) {
	if userName == "" {
		u, err := user.Current()
		if err == nil {
			userName = u.Username
		}
	}
	return &Session{
		Environment: env,
		Path:        path,
		FromVersion: from,
		ToVersion:   to,
		User:        userName,
	}, nil
}

// BuildEntry constructs an audit Entry from the session and a slice of diff changes.
func (s *Session) BuildEntry(changes []diff.Change) Entry {
	var summary diff.Summary
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			summary.Added++
		case diff.Removed:
			summary.Removed++
		case diff.Modified:
			summary.Modified++
		case diff.Unchanged:
			summary.Unchanged++
		}
	}
	return Entry{
		Timestamp:   time.Now().UTC(),
		Environment: s.Environment,
		Path:        s.Path,
		FromVersion: s.FromVersion,
		ToVersion:   s.ToVersion,
		User:        s.User,
		Changes:     changes,
		Summary:     summary,
	}
}
