// Package scaffold verifies and generates workspace files:
// .claude/skills/, .claude/agents/, opencode.json instructions and CLAUDE.md.
package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vantrilex/internal/doctor"
)

// Toolkit is one of the 8 designated skills/agent collections provisioned
// into every new project workspace.
type Toolkit struct {
	Slug      string // bridge filename slug, e.g. "vercel-skills"
	Name      string // display name, e.g. "vercel-labs/skills"
	Source    string // clone URL (all 8 are git-backed)
	Dir       string // expected cache dir under $HOME/.vantrilex/toolkit
	SkillHint string // where the key skill lives inside the cache
}

// toolkitMeta carries display names and skill hints keyed by cache dir.
var toolkitMeta = map[string]struct{ name, hint string }{
	"skills":                           {"mattpocock/skills", "TypeScript and agent skills; glob **/SKILL.md."},
	"ponytail":                         {"dietrichgebert/ponytail", "Multi-agent orchestration; glob **/SKILL.md."},
	"guard-skills":                     {"amElnagdy/guard-skills", "Guardrails and safety skills; glob **/SKILL.md."},
	"Anthropic-Cybersecurity-Skills":   {"mukul975/Anthropic-Cybersecurity-Skills", "Cybersecurity skills; glob **/SKILL.md."},
	"agency-agents":                    {"msitarzewski/agency-agents", "Agency sub-agents; glob **/SKILL.md and agents."},
	"everything-claude-code":           {"worldflowai/everything-claude-code", "Claude Code setups and workflows; glob **/SKILL.md."},
	"vercel-skills":                    {"vercel-labs/skills", "Key skill: find-skills at skills/find-skills/SKILL.md."},
	"anthropic-skills":                 {"anthropics/skills (official)", "Key skill: skill-creator at skills/skill-creator/SKILL.md."},
}

// Toolkits returns all 8 designated toolkits, sourced from doctor.ToolkitRepos
// (single source of truth for the clone set).
func Toolkits() []Toolkit {
	out := make([]Toolkit, 0, len(doctor.ToolkitRepos))
	for _, r := range doctor.ToolkitRepos {
		meta, ok := toolkitMeta[r.Dir]
		name, hint := r.Dir, "glob **/SKILL.md"
		if ok {
			name, hint = meta.name, meta.hint
		}
		out = append(out, Toolkit{
			Slug:      slugify(r.Dir),
			Name:      name,
			Source:    r.URL,
			Dir:       r.Dir,
			SkillHint: hint,
		})
	}
	return out
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	for strings.Contains(out, "--") {
		out = strings.ReplaceAll(out, "--", "-")
	}
	return out
}

// BridgeRelPath returns the workspace-relative bridge file for a toolkit.
func BridgeRelPath(t Toolkit) string {
	return ".claude/skills/toolkit-" + t.Slug + ".md"
}

// Item is one scaffold check/file.
type Item struct {
	Path    string // relative to workspace
	Kind    string // "dir" | "file"
	Present bool
	Note    string
}

// Plan returns the expected workspace layout: base files plus one bridge
// file per designated toolkit (all 8).
func Plan() []Item {
	plan := []Item{
		{Path: ".claude/skills", Kind: "dir", Note: "Workflow skills live here"},
		{Path: ".claude/agents", Kind: "dir", Note: "Sub-agent definitions"},
		{Path: "opencode.json", Kind: "file", Note: "OpenCode instructions + providers"},
		{Path: "CLAUDE.md", Kind: "file", Note: "Project instructions for Claude"},
	}
	for _, t := range Toolkits() {
		plan = append(plan, Item{Path: BridgeRelPath(t), Kind: "file", Note: "Toolkit bridge: " + t.Name})
	}
	return plan
}

// Check verifies which items already exist in workspace.
func Check(workspace string) []Item {
	plan := Plan()
	for i, it := range plan {
		_, err := os.Stat(filepath.Join(workspace, filepath.FromSlash(it.Path)))
		plan[i].Present = err == nil
	}
	return plan
}

// AllPresent reports whether scaffolding is complete.
func AllPresent(workspace string) bool {
	for _, it := range Check(workspace) {
		if !it.Present {
			return false
		}
	}
	return true
}

const claudeMdTemplate = `# Project Instructions (Vantrilex)

This workspace was provisioned by the Vantrilex launcher.

## Runners
- Claude Code CLI (` + "`claude`" + `): streaming tools, sub-agents.
- OpenCode Agent (` + "`opencode`" + `): multi-provider routing + Zen.
- OpenAI Codex CLI (` + "`codex`" + `): GPT-5.6 execution.

## Skills
Workflow skills live in ` + "`.claude/skills/`" + `, sub-agents in ` + "`.claude/agents/`" + `.
Shared toolkit cache: ` + "`$HOME/.vantrilex/toolkit`" + `.
Per-toolkit bridges: ` + "`.claude/skills/toolkit-*.md`" + ` (one for each of the 8 toolkits).

## External skills
- vercel-labs/skills@find-skills (cache: $HOME/.vantrilex/toolkit/vercel-skills)
- Anthropic Official skill-creator (cache: $HOME/.vantrilex/toolkit/anthropic-skills)

## Rules
- Keep responses strictly in English.
- Prefer small, reversible changes with tests.
`

// toolkitsMarkdown renders the 8-toolkit reference section for CLAUDE.md.
func toolkitsMarkdown() string {
	var b strings.Builder
	b.WriteString("\n## Toolkits (8 provisioned)\n\n")
	b.WriteString("Each toolkit is cached under `$HOME/.vantrilex/toolkit/` and bridged\n")
	b.WriteString("into `.claude/skills/toolkit-<slug>.md`. Load a bridge, then glob the\n")
	b.WriteString("cache for `**/SKILL.md` to activate skills.\n\n")
	for _, t := range Toolkits() {
		fmt.Fprintf(&b, "- %s — `%s` (cache `%s`): %s\n", t.Name, t.Source, t.Dir, t.SkillHint)
	}
	return b.String()
}

const opencodeJSONTemplate = `{
  "$schema": "https://opencode.ai/config.json",
  "instructions": ["CLAUDE.md", ".claude/skills/*.md", ".claude/skills/**/SKILL.md", ".claude/agents/*.md"],
  "model": "opencode/zen"
}
`

const starterSkillTemplate = `# %s

Starter skill scaffolded by Vantrilex.

## Workflow
1. Read the task and relevant repo files first.
2. Plan the smallest reversible change.
3. Implement, then verify (build/tests).
4. Summarize in English.

External references:
- vercel-labs/skills@find-skills
- Anthropic Official skill-creator
`

const starterAgentTemplate = `---
name: %s
description: Starter sub-agent scaffolded by Vantrilex.
---

You are the %s sub-agent. Work in English, keep changes small and verified.
`

// Apply creates missing dirs/files (non-destructive: never overwrites).
func Apply(workspace, runnerID, modelID, effort string) ([]string, error) {
	var created []string
	mk := func(rel string) (string, error) {
		abs := filepath.Join(workspace, filepath.FromSlash(rel))
		if err := os.MkdirAll(abs, 0o755); err != nil {
			return "", err
		}
		return abs, nil
	}
	writeOnce := func(rel, content string) error {
		abs := filepath.Join(workspace, filepath.FromSlash(rel))
		if _, err := os.Stat(abs); err == nil {
			return nil // keep existing
		}
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			return err
		}
		created = append(created, rel)