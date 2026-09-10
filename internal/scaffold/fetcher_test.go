package scaffold

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vantrilex/internal/catalog"
)

func TestFetchBodyLive(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("# skill body"))
	}))
	defer srv.Close()
	body, err := FetchBody(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "# skill body" {
		t.Fatalf("body=%q", body)
	}
}

func TestApplySelectionsOfflineFallback(t *testing.T) {
	ws := filepath.Join(t.TempDir(), "proj")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}
	// Point RawURLs at unroutable address to force fallback templates.
	bad := "http://127.0.0.1:1/nope"
	agents := []catalog.RegistryItem{{Name: "Lead System Architect", Desc: "orchestrator", Source: "Architecture", RawURL: bad}}
	skills := []catalog.RegistryItem{{Name: "find-skills", Desc: "discover", Source: "vercel-labs/skills", RawURL: bad}}
	plugins := []catalog.RegistryItem{{Name: "commit-commands", Desc: "commits", Install: "official/commit-commands", RawURL: bad}}
	hooks := []catalog.RegistryItem{{Name: "format-on-edit", Desc: "format", Source: "PostToolUse", Install: "PostToolUse:format-on-edit", RawURL: bad}}
	mcps := []catalog.RegistryItem{{Name: "fetch", Desc: "web fetch", Install: "npx -y @modelcontextprotocol/server-fetch", RawURL: bad}}
	created, err := ApplySelections(context.Background(), ws, agents, skills, plugins, hooks, mcps, "opencode", "opencode/zen", "medium")
	if err != nil {
		t.Fatal(err)
	}
	if len(created) == 0 {
		t.Fatal("expected created files")
	}
	for _, want := range []string{".claude/skills/find-skills/SKILL.md", ".claude/agents/lead-system-architect.md", "opencode.json", ".mcp.json", "CLAUDE.md", ".claude/settings.json"} {
		if _, err := os.Stat(filepath.Join(ws, filepath.FromSlash(want))); err != nil {
			t.Fatalf("missing %s", want)
		}
	}
	oj, _ := os.ReadFile(filepath.Join(ws, "opencode.json"))
	if !strings.Contains(string(oj), "fetch") {
		t.Fatal("opencode.json missing mcp key")
	}
	cm, _ := os.ReadFile(filepath.Join(ws, "CLAUDE.md"))
	if !strings.Contains(string(cm), "Lead System Architect") {
		t.Fatal("manifest missing agent")
	}
}
