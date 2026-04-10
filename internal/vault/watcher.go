package vault

import (
	"context"
	"time"
)

// WatcherConfig holds configuration for polling a Vault secret path.
type WatcherConfig struct {
	Path     SecretPath
	Interval time.Duration
	MaxDrift time.Duration
}

// WatchEvent is emitted each time the watcher polls a secret.
type WatchEvent struct {
	Path      SecretPath
	Version   int
	ChangedAt time.Time
	Err       error
}

// Watcher polls a Vault KV secret path at a fixed interval and emits
// WatchEvents on a channel whenever the version changes.
type Watcher struct {
	cfg    WatcherConfig
	reader SecretReader
}

// SecretReader is a minimal interface consumed by Watcher.
type SecretReader interface {
	ReadSecret(ctx context.Context, path SecretPath) (*SecretVersion, error)
}

// NewWatcher creates a Watcher with the provided config and reader.
// An error is returned if the path is invalid or the interval is zero.
func NewWatcher(cfg WatcherConfig, reader SecretReader) (*Watcher, error) {
	if err := cfg.Path.Validate(); err != nil {
		return nil, err
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 30 * time.Second
	}
	if cfg.MaxDrift <= 0 {
		cfg.MaxDrift = cfg.Interval / 2
	}
	return &Watcher{cfg: cfg, reader: reader}, nil
}

// Watch starts polling and sends WatchEvents to the returned channel.
// The channel is closed when ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) <-chan WatchEvent {
	ch := make(chan WatchEvent, 1)
	go func() {
		defer close(ch)
		var lastVersion int
		ticker := time.NewTicker(w.cfg.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				sv, err := w.reader.ReadSecret(ctx, w.cfg.Path)
				if err != nil {
					ch <- WatchEvent{Path: w.cfg.Path, Err: err}
					continue
				}
				if sv.Version != lastVersion {
					lastVersion = sv.Version
					ch <- WatchEvent{
						Path:      w.cfg.Path,
						Version:   sv.Version,
						ChangedAt: t,
					}
				}
			}
		}
	}()
	return ch
}
