// Package doctor implements the Preflight Dependency Doctor: it scans the
// host for required runtime tools, installs missing ones via winget/npm,
// and shallow-clones the agent/skills toolkit in parallel.
package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// ToolkitRepos are shallow-cloned into $HOME/.vantrilex/toolkit.
// The full set covers all 8 designated toolkits/skills: the 6 agent/skills
// collections plus the vercel find-skills repo and the Anthropic official
// skills repo (home of skill-creator).
var ToolkitRepos = []struct {
	URL string
	Dir string
}{
	{"https://github.com/mattpocock/skills", "skills"},
	{"https://github.com/dietrichgebert/ponytail", "ponytail"},
	{"https://github.com/amElnagdy/guard-skills", "guard-skills"},
	{"https://github.com/mukul975/Anthropic-Cybersecurity-Skills", "Anthropic-Cybersecurity-Skills"},
	{"https://github.com/msitarzewski/agency-agents", "agency-agents"},
	{"https://github.com/worldflowai/everything-claude-code", "everything-claude-code"},
	{"https://github.com/vercel-labs/skills", "vercel-skills"},
	{"https://github.com/anthropics/skills", "anthropic-skills"},
}

// ExternalSkills are the named skills resolved inside the toolkit clones:
// find-skills lives in vercel-labs/skills, skill-creator in anthropics/skills.
var ExternalSkills = []string{
	"vercel-labs/skills@find-skills",
	"anthropic/skill-creator (Official skill-creator)",
}

// Dep describes one prerequisite.
type Dep struct {
	Key        string // stable id: git, node, npx, chafa, claude, opencode, codex
	Label      string
	Binary     string   // binary probed on PATH
	Fallbacks  []string // alternate binaries (e.g. nodejs)
	VersionArg string
	Kind       string // "binary" | "npm"
	NpmPkg     string // for npm kinds
	WingetID   string // winget package id (windows)
	Optional   bool
}

// Deps is the full prerequisite list.
func Deps() []Dep {
	return []Dep{
		{Key: "git", Label: "Git (repo sync)", Binary: "git", VersionArg: "--version", WingetID: "Git.Git"},
		{Key: "node", Label: "Node.js (runners)", Binary: "node", VersionArg: "--version", WingetID: "OpenJS.NodeJS.LTS"},
		{Key: "npx", Label: "npx (runner exec)", Binary: "npx", VersionArg: "--version"},
		{Key: "chafa", Label: "Chafa (logo render)", Binary: "chafa", VersionArg: "--version", WingetID: "HansPetterJansson.Chafa", Optional: true},
		{Key: "claude", Label: "Claude Code CLI", Binary: "claude", VersionArg: "--version", Kind: "npm", NpmPkg: "@anthropic-ai/claude-code"},
		{Key: "opencode", Label: "OpenCode Agent", Binary: "opencode", Fallbacks: []string{"opencode-ai"}, VersionArg: "--version", Kind: "npm", NpmPkg: "opencode-ai"},
		{Key: "codex", Label: "OpenAI Codex CLI", Binary: "codex", VersionArg: "--version", Kind: "npm", NpmPkg: "@openai/codex", Optional: true},
	}
}

// Status is a check result for one dep.
type Status struct {