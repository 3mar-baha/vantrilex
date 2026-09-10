package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplyAndCheck(t *testing.T) {
	ws := filepath.Join(t.TempDir(), "proj")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}
	created, err := Apply(ws, "opencode", "opencode/zen", "medium")
	if err != nil {
		t.Fatal(err)
	}
	if len(created) == 0 {
		t.Fatal("expected created files on fresh workspace")
	}
	if !AllPresent(ws) {
		t.Fatal("expected all scaffold items present after Apply")
	}
	// Second apply must be non-destructive (no overwrite, no error).
	if _, err := Apply(ws, "opencode", "opencode/zen", "medium"); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"CLAUDE.md", "opencode.json", ".claude/skills/starter/SKILL.md", ".claude/agents/planner.md"} {
		if _, err := os.Stat(filepath.Join(ws, filepath.FromSlash(want))); err != nil {
			t.Fatalf("missing %s: %v", want, err)
		}
	}
	// All 8 toolkit bridges must be injected.
	if got := Toolkits(); len(got) != 8 {
		t.Fatalf("expected 8 toolkits, got %d", len(got))
	}
	for _, tk := range Toolkits() {
		rel := BridgeRelPath(tk)
		data, err := os.ReadFile(filepath.Join(ws, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("missing toolkit bridge %s: %v", rel, err)
		}
		if !strings.Contains(string(data), tk.Source) {
			t.Fatalf("bridge %s does not reference source %s", rel, tk.Source)
		}
	}
	// opencode.json must pick up the bridge files.
	oj, err := os.ReadFile(filepath.Join(ws, "opencode.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(oj), ".claude/skills/*.md") {
		t.Fatal("opencode.json instructions do not include .claude/skills/*.md bridges")
	}
	// CLAUDE.md must document all 8 toolkits.
	cm, err := os.ReadFile(filepath.Join(ws, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tk := range Toolkits() {
		if !strings.Contains(string(cm), tk.Source) {
			t.Fatalf("CLAUDE.md does not document toolkit %s", tk.Source)
		}
	}
}
