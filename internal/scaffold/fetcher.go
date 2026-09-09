// On-demand asset fetcher: zero bulk clones. Selected registry items are
// fetched over HTTPS only when needed and written straight into the target
// project workspace. Offline-safe: starter bodies are embedded so Apply
// succeeds without network; HTTPS upgrades content best-effort.
package scaffold

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vantrilex/internal/catalog"
)

var fetchHTTP = &http.Client{Timeout: 20 * time.Second}

// FetchBody downloads a raw definition over HTTPS (single retry).
func FetchBody(ctx context.Context, url string) ([]byte, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil, fmt.Errorf("empty url")
	}
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := fetchHTTP.Do(req)
		if err != nil {
			last = err
			continue
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		resp.Body.Close()
		if err != nil {
			last = err
			continue
		}
		if resp.StatusCode != 200 {
			last = fmt.Errorf("http %d for %s", resp.StatusCode, url)
			continue
		}
		if len(body) == 0 {
			last = fmt.Errorf("empty body for %s", url)
			continue
		}
		return body, nil
	}
	return nil, last
}

func writeOnceFile(abs, content string, created *[]string, rel string) error {
	if _, err := os.Stat(abs); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		return err
	}
	*created = append(*created, rel)
	return nil
}

func sanitize(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		case r == ' ' || r == '/' || r == '.':
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	for strings.Contains(out, "--") {
		out = strings.ReplaceAll(out, "--", "-")
	}
	if out == "" {
		out = "asset"
	}
	return out
}

// ApplySelections provisions only the checked items into the workspace.
func ApplySelections(ctx context.Context, workspace string, agents, skills, plugins, hooks, mcps []catalog.RegistryItem, runnerID, modelID, effort string) ([]string, error) {
	var created []string
	base := []string{".claude/skills", ".claude/agents", ".claude/plugins", ".claude/hooks"}
	for _, d := range base {
		if err := os.MkdirAll(filepath.Join(workspace, filepath.FromSlash(d)), 0o755); err != nil {
			return created, err
		}
	}
	// Skills: try HTTPS raw, fall back to starter template.
	for _, s := range skills {
		rel := ".claude/skills/" + sanitize(s.Name) + "/SKILL.md"
		abs := filepath.Join(workspace, filepath.FromSlash(rel))
		if _, err := os.Stat(abs); err == nil {
			continue
		}
		content := fmt.Sprintf(starterSkillTemplate, s.Name) + "\nSource: " + s.Source + "\n" + s.Desc + "\n"
		if body, err := FetchBody(ctx, s.RawURL); err == nil {
			content = "# " + s.Name + "\n\n" + string(body)
		}
		if err := writeOnceFile(abs, content, &created, rel); err != nil {
			return created, err
		}
	}
	// Agents.
	for _, a := range agents {
		rel := ".claude/agents/" + sanitize(a.Name) + ".md"
		abs := filepath.Join(workspace, filepath.FromSlash(rel))
		if _, err := os.Stat(abs); err == nil {
			continue
		}
		content := fmt.Sprintf(starterAgentTemplate, a.Name, a.Name) + "\nRole: " + a.Source + "\n" + a.Desc + "\n"
		if body, err := FetchBody(ctx, a.RawURL); err == nil {
			content = string(body)
		}
		if err := writeOnceFile(abs, content, &created, rel); err != nil {
			return created, err
		}
	}
	// Plugins: record install refs.
	if len(plugins) > 0 {
		rel := ".claude/plugins/manifest.json"
		abs := filepath.Join(workspace, filepath.FromSlash(rel))
		if _, err := os.Stat(abs); err != nil {
			refs := []string{}
			for _, p := range plugins {
				refs = append(refs, p.Install)
			}
			b, _ := json.MarshalIndent(map[string]any{"plugins": refs}, "", "  ")
			if err := writeOnceFile(abs, string(b)+"\n", &created, rel); err != nil {
				return created, err
			}
		}
	}
	// Hooks: map events in settings.json.
	if len(hooks) > 0 {
		rel := ".claude/settings.json"
		abs := filepath.Join(workspace, filepath.FromSlash(rel))
		var cur map[string]any
		if data, err := os.ReadFile(abs); err == nil {
			_ = json.Unmarshal(data, &cur)
		}
		if cur == nil {
			cur = map[string]any{}
		}
		hmap, _ := cur["hooks"].(map[string]any)
		if hmap == nil {
			hmap = map[string]any{}
		}
		changed := false
		for _, h := range hooks {
			ev := h.Source
			if ev == "" {
				ev = "PostToolUse"
			}
			key := ev + ":" + h.Name
			if _, ok := hmap[key]; !ok {
				hmap[key] = map[string]any{"command": h.Install, "event": ev}
				changed = true
			}
			// Best-effort fetch of hook body into hooks dir.
			hrel := ".claude/hooks/" + sanitize(h.Name) + ".json"
			habs := filepath.Join(workspace, filepath.FromSlash(hrel))
			if _, err := os.Stat(habs); err != nil {
				hcontent := fmt.Sprintf("{\"name\": %q, \"event\": %q, \"note\": %q}", h.Name, ev, h.Desc)
				if body, err := FetchBody(ctx, h.RawURL); err == nil {
					hcontent = string(body)
				}
				_ = writeOnceFile(habs, hcontent, &created, hrel)
			}
		}
		if changed {
			cur["hooks"] = hmap
			b, _ := json.MarshalIndent(cur, "", "  ")
			if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
				return created, err
			}
			if _, err := os.Stat(abs); err != nil {
				if err := os.WriteFile(abs, append(b, '\n'), 0o644); err != nil {
					return created, err
				}
				created = append(created, rel)
			}
		}
	}
	// MCP servers: inject into opencode.json + .mcp.json.
	if len(mcps) > 0 {
		if err := injectMCP(workspace, mcps, &created); err != nil {
			return created, err
		}
	}
	// Manifest into CLAUDE.md (append active component list).
	if err := appendManifest(workspace, agents, skills, plugins, hooks, mcps, runnerID, modelID, effort, &created); err != nil {
		return created, err
	}
	return created, nil
}

func injectMCP(workspace string, mcps []catalog.RegistryItem, created *[]string) error {
	// opencode.json
	ojPath := filepath.Join(workspace, "opencode.json")
	var oj map[string]any
	if data, err := os.ReadFile(ojPath); err == nil {
		_ = json.Unmarshal(data, &oj)
	}
	if oj == nil {
		oj = map[string]any{}
	}
	mcpMap, _ := oj["mcp"].(map[string]any)
	if mcpMap == nil {
		mcpMap = map[string]any{}
	}
	for _, m := range mcps {
		key := sanitize(m.Name)
		if _, ok := mcpMap[key]; !ok {
			parts := strings.Fields(m.Install)
			cmd := "npx"
			args := []string{}
			if len(parts) > 0 {
				cmd = parts[0]
				args = parts[1:]
			}
			mcpMap[key] = map[string]any{"command": cmd, "args": args, "type": "local", "enabled": true}
		}
	}
	oj["mcp"] = mcpMap
	b, _ := json.MarshalIndent(oj, "", "  ")
	if err := os.WriteFile(ojPath, append(b, '\n'), 0o644); err != nil {
		return err
	}
	// .mcp.json
	mjPath := filepath.Join(workspace, ".mcp.json")
	var mj map[string]any
	if data, err := os.ReadFile(mjPath); err == nil {
		_ = json.Unmarshal(data, &mj)