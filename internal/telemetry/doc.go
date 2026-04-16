// Package telemetry provides lightweight span-based timing for vaultpipe
// operations. It is intentionally minimal — no external dependencies, no
// network export — and is designed for surfacing latency information to the
// developer at the end of a CLI run via --timings.
package telemetry
