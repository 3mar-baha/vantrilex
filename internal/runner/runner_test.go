package runner

import "testing"

import "vantrilex/internal/catalog"

func TestBuildCommand(t *testing.T) {
	r, _ := catalog.RunnerByID(catalog.RunnerClaude)
	name, args := BuildCommand(r, "anthropic/claude-sonnet-4-5", "medium", "C:\\tmp\\x")
	if name != "claude" || len(args) != 2 || args[0] != "--model" {
		t.Fatalf("unexpected claude argv: %s %v", name, args)
	}
	r2, _ := catalog.RunnerByID(catalog.RunnerCodex)
	name2, args2 := BuildCommand(r2, "openai/gpt-5.6", "max", "C:\\tmp\\x")
	if name2 != "codex" || len(args2) != 4 {
		t.Fatalf("unexpected codex argv: %s %v", name2, args2)
	}
	if got := TruncPath("C:\\very\\long\\path\\to\\project\\folder\\name", 20); len(got) > 23 {
		t.Fatalf("truncation failed: %q", got)
	}
}
