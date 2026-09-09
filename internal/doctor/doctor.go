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
	Dep     Dep
	Found   bool
	Version string
	Err     string
}

// LookPath checks PATH for binary or fallbacks.
func LookPath(bin string, fallbacks []string) (string, bool) {
	if p, err := exec.LookPath(bin); err == nil {
		return p, true
	}
	for _, f := range fallbacks {
		if p, err := exec.LookPath(f); err == nil {
			return p, true
		}
	}
	return "", false
}

// CheckOne probes a single dependency.
func CheckOne(d Dep) Status {
	st := Status{Dep: d}
	path, ok := LookPath(d.Binary, d.Fallbacks)
	if !ok {
		st.Found = false
		st.Err = "not on PATH"
		return st
	}
	st.Found = true
	cmd := exec.Command(path, d.VersionArg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		st.Version = "installed"
		return st
	}
	line := strings.TrimSpace(strings.SplitN(string(out), "\n", 2)[0])
	if len(line) > 64 {
		line = line[:64]
	}
	st.Version = line
	return st
}

// CheckAll scans every prerequisite.
func CheckAll() []Status {
	deps := Deps()
	out := make([]Status, 0, len(deps))
	for _, d := range deps {
		out = append(out, CheckOne(d))
	}
	return out
}

// MissingRequired returns statuses that are absent and not optional.
func MissingRequired(all []Status) []Status {
	var out []Status
	for _, s := range all {
		if !s.Found && !s.Dep.Optional {
			out = append(out, s)
		}
	}
	return out
}

// MissingAny returns all absent statuses including optional.
func MissingAny(all []Status) []Status {
	var out []Status
	for _, s := range all {