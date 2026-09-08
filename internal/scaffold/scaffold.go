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