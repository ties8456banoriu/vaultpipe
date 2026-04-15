// Package refresh provides functionality for watching and automatically
// re-injecting secrets from Vault when they are close to expiry or on demand.
package refresh

import (
	"context"
	"log"
	"time"
)

// SecretFetcher abstracts fetching secrets from a backend (e.g. Vault).
type SecretFetcher interface {
	GetSecret(path string) (map[string]string, error)
}

// EnvWriter abstracts writing secrets to an env file.
type EnvWriter interface {
	Write(secrets map[string]string) error
}

// Watcher polls Vault at a given interval and rewrites the env file when
// secrets change.
type Watcher struct {
	fetcher  SecretFetcher
	writer   EnvWriter
	path     string
	interval time.Duration
	logger   *log.Logger
}

// NewWatcher creates a Watcher that polls the given secret path every interval.
func NewWatcher(fetcher SecretFetcher, writer EnvWriter, path string, interval time.Duration, logger *log.Logger) *Watcher {
	return &Watcher{
		fetcher:  fetcher,
		writer:   writer,
		path:     path,
		interval: interval,
		logger:   logger,
	}
}

// Start begins the polling loop. It blocks until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) error {
	w.logger.Printf("refresh: starting watcher for path=%s interval=%s", w.path, w.interval)

	if err := w.refresh(); err != nil {
		w.logger.Printf("refresh: initial fetch failed: %v", err)
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.refresh(); err != nil {
				w.logger.Printf("refresh: fetch failed: %v", err)
			}
		case <-ctx.Done():
			w.logger.Println("refresh: watcher stopped")
			return ctx.Err()
		}
	}
}

func (w *Watcher) refresh() error {
	secrets, err := w.fetcher.GetSecret(w.path)
	if err != nil {
		return err
	}
	return w.writer.Write(secrets)
}
