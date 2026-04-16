// Package hooks provides pre/post secret-fetch lifecycle hooks
// that can run shell commands or log events at key pipeline stages.
package hooks

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Stage represents a lifecycle point in the pipeline.
type Stage string

const (
	StagePre  Stage = "pre"
	StagePost Stage = "post"
)

// Hook defines a single lifecycle hook.
type Hook struct {
	Stage   Stage
	Command string
}

// Runner executes registered hooks for a given stage.
type Runner struct {
	hooks []Hook
}

// NewRunner creates a Runner with the provided hooks.
func NewRunner(hooks []Hook) *Runner {
	return &Runner{hooks: hooks}
}

// Run executes all hooks matching the given stage.
// Each hook command is run via the shell. Errors are collected and returned.
func (r *Runner) Run(ctx context.Context, stage Stage) error {
	var errs []string
	for _, h := range r.hooks {
		if h.Stage != stage {
			continue
		}
		if err := runCommand(ctx, h.Command); err != nil {
			errs = append(errs, fmt.Sprintf("hook %q: %v", h.Command, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("hook errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

func runCommand(ctx context.Context, command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	return cmd.Run()
}

// ParseRules parses hook definitions from string slice of "stage:command" pairs.
func ParseRules(rules []string) ([]Hook, error) {
	var hooks []Hook
	for _, r := range rules {
		idx := strings.Index(r, ":")
		if idx < 0 {
			return nil, fmt.Errorf("invalid hook rule %q: expected stage:command", r)
		}
		stage := Stage(strings.TrimSpace(r[:idx]))
		if stage != StagePre && stage != StagePost {
			return nil, fmt.Errorf("unknown stage %q: must be pre or post", stage)
		}
		cmd := strings.TrimSpace(r[idx+1:])
		if cmd == "" {
			return nil, fmt.Errorf("hook rule %q has empty command", r)
		}
		hooks = append(hooks, Hook{Stage: stage, Command: cmd})
	}
	return hooks, nil
}
