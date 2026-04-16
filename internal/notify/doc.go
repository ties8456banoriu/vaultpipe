// Package notify emits notifications when vaultpipe successfully injects
// secrets into the local environment. Notifications are written to a
// configurable io.Writer and, optionally, forwarded to the host OS
// notification system (macOS: osascript, Linux: notify-send).
package notify
