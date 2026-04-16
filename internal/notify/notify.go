// Package notify provides post-write desktop and log notifications
// when secrets are successfully injected into the environment.
package notify

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"time"
)

// Event holds metadata about a completed secret injection.
type Event struct {
	Profile   string
	KeyCount  int
	OccurredAt time.Time
}

// Notifier sends notifications about injection events.
type Notifier struct {
	out io.Writer
	desktop bool
}

// New creates a Notifier. If desktop is true it attempts OS-level
// notifications in addition to writing to out.
func New(out io.Writer, desktop bool) *Notifier {
	return &Notifier{out: out, desktop: desktop}
}

// Notify emits a notification for the given event.
func (n *Notifier) Notify(e Event) error {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}
	msg := fmt.Sprintf("[vaultpipe] %s — injected %d secret(s) at %s",
		e.Profile, e.KeyCount, e.OccurredAt.Format(time.RFC3339))

	if _, err := fmt.Fprintln(n.out, msg); err != nil {
		return fmt.Errorf("notify: write: %w", err)
	}

	if n.desktop {
		_ = sendDesktop(e.Profile, msg)
	}
	return nil
}

// sendDesktop attempts a best-effort OS notification; errors are ignored.
func sendDesktop(title, body string) error {
	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q`, body, title)
		return exec.Command("osascript", "-e", script).Run()
	case "linux":
		return exec.Command("notify-send", title, body).Run()
	default:
		return nil
	}
}
