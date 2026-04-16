package hooks_test

import (
	"context"
	"testing"

	"github.com/yourusername/vaultpipe/internal/hooks"
)

func TestParseRules_Valid(t *testing.T) {
	rules := []string{"pre:echo hello", "post:echo bye"}
	h, err := hooks.ParseRules(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h) != 2 {
		t.Fatalf("expected 2 hooks, got %d", len(h))
	}
	if h[0].Stage != hooks.StagePre {
		t.Errorf("expected pre, got %s", h[0].Stage)
	}
	if h[1].Stage != hooks.StagePost {
		t.Errorf("expected post, got %s", h[1].Stage)
	}
}

func TestParseRules_MissingColon(t *testing.T) {
	_, err := hooks.ParseRules([]string{"preecho"})
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseRules_UnknownStage(t *testing.T) {
	_, err := hooks.ParseRules([]string{"during:echo hi"})
	if err == nil {
		t.Fatal("expected error for unknown stage")
	}
}

func TestParseRules_EmptyCommand(t *testing.T) {
	_, err := hooks.ParseRules([]string{"pre:"})
	if err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestRun_PreHook_Success(t *testing.T) {
	h := []hooks.Hook{{Stage: hooks.StagePre, Command: "echo hello"}}
	r := hooks.NewRunner(h)
	if err := r.Run(context.Background(), hooks.StagePre); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_PostHook_NotTriggeredForPre(t *testing.T) {
	h := []hooks.Hook{{Stage: hooks.StagePost, Command: "echo bye"}}
	r := hooks.NewRunner(h)
	if err := r.Run(context.Background(), hooks.StagePre); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_FailingCommand_ReturnsError(t *testing.T) {
	h := []hooks.Hook{{Stage: hooks.StagePre, Command: "false"}}
	r := hooks.NewRunner(h)
	if err := r.Run(context.Background(), hooks.StagePre); err == nil {
		t.Fatal("expected error from failing command")
	}
}

func TestRun_EmptyHooks_NoError(t *testing.T) {
	r := hooks.NewRunner(nil)
	if err := r.Run(context.Background(), hooks.StagePre); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
