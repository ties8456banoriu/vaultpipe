// Package audit provides structured, append-only audit logging for vaultpipe.
//
// Each secret fetch and .env write operation is recorded as a JSON-lines entry
// containing a timestamp, event type, and relevant metadata (secret path,
// output file, key count). This makes it straightforward to pipe audit output
// to a file or log aggregation system:
//
//	vaultpipe run 2>> audit.log
//
// Events are written to stderr by default so they remain separate from normal
// program output, but any io.Writer can be supplied via NewLogger.
//
// Supported event types:
//
//	- "secret_fetch"  – a secret was retrieved from Vault
//	- "env_write"     – a .env file was written with fetched secrets
//	- "secret_error"  – an error occurred while fetching a secret
package audit
