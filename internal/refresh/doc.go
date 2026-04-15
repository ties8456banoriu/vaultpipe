// Package refresh implements an automatic secret refresh mechanism for vaultpipe.
//
// A Watcher polls HashiCorp Vault at a configurable interval and rewrites the
// local .env file whenever new secret values are retrieved. This allows
// long-running development processes to pick up rotated credentials without
// manual intervention.
//
// Basic usage:
//
//	w := refresh.NewWatcher(vaultClient, envWriter, "secret/data/myapp", 5*time.Minute, logger)
//	if err := w.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package refresh
