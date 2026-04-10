package vault

import (
	"context"
	"errors"
	"testing"
	"time"
)

// stubReader is a test double for SecretReader.
type stubReader struct {
	versions []int
	calls    int
	err      error
}

func (s *stubReader) ReadSecret(_ context.Context, _ SecretPath) (*SecretVersion, error) {
	if s.err != nil {
		return nil, s.err
	}
	v := s.versions[s.calls%len(s.versions)]
	s.calls++
	return &SecretVersion{Version: v, Data: map[string]interface{}{}}, nil
}

func validPath(t *testing.T) SecretPath {
	t.Helper()
	p, err := NewSecretPath("secret", "myapp/config")
	if err != nil {
		t.Fatalf("NewSecretPath: %v", err)
	}
	return p
}

func TestNewWatcher_Valid(t *testing.T) {
	w, err := NewWatcher(WatcherConfig{
		Path:     validPath(t),
		Interval: time.Second,
	}, &stubReader{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestNewWatcher_DefaultInterval(t *testing.T) {
	w, err := NewWatcher(WatcherConfig{Path: validPath(t)}, &stubReader{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.cfg.Interval != 30*time.Second {
		t.Errorf("expected default interval 30s, got %v", w.cfg.Interval)
	}
}

func TestNewWatcher_InvalidPath(t *testing.T) {
	_, err := NewWatcher(WatcherConfig{}, &stubReader{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestWatcher_EmitsOnVersionChange(t *testing.T) {
	reader := &stubReader{versions: []int{1, 2}}
	w, _ := NewWatcher(WatcherConfig{
		Path:     validPath(t),
		Interval: 10 * time.Millisecond,
	}, reader)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ch := w.Watch(ctx)
	var events []WatchEvent
	for e := range ch {
		events = append(events, e)
	}
	if len(events) < 2 {
		t.Errorf("expected at least 2 events, got %d", len(events))
	}
}

func TestWatcher_EmitsErrorEvent(t *testing.T) {
	reader := &stubReader{err: errors.New("vault unavailable")}
	w, _ := NewWatcher(WatcherConfig{
		Path:     validPath(t),
		Interval: 10 * time.Millisecond,
	}, reader)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	ch := w.Watch(ctx)
	var errCount int
	for e := range ch {
		if e.Err != nil {
			errCount++
		}
	}
	if errCount == 0 {
		t.Error("expected at least one error event")
	}
}
