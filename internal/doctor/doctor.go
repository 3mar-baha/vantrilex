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
		if !s.Found {
			out = append(out, s)
		}
	}
	return out
}

// InstallCommands returns the human-readable install plan for a dep.
func InstallCommands(d Dep) []string {
	var cmds []string
	if runtime.GOOS == "windows" && d.WingetID != "" {
		cmds = append(cmds, fmt.Sprintf("winget install -e --id %s", d.WingetID))
	}
	if d.Kind == "npm" && d.NpmPkg != "" {
		cmds = append(cmds, fmt.Sprintf("npm install -g %s", d.NpmPkg))
	}
	if len(cmds) == 0 && runtime.GOOS != "windows" {
		cmds = append(cmds, fmt.Sprintf("install %s via system package manager", d.Binary))
	}
	return cmds
}

// InstallOne attempts a non-destructive install of a missing dep.
// It prefers winget on Windows, then npm for npm-kind deps.
func InstallOne(d Dep, log func(string)) error {
	if _, ok := LookPath(d.Binary, d.Fallbacks); ok {
		return nil // already present (race-safe)
	}
	if runtime.GOOS == "windows" && d.WingetID != "" {
		log(fmt.Sprintf("winget: installing %s ...", d.WingetID))
		cmd := exec.Command("winget", "install", "-e", "--silent", "--accept-package-agreements", "--accept-source-agreements", "--id", d.WingetID)
		out, err := cmd.CombinedOutput()
		if err == nil {
			return nil
		}
		log(fmt.Sprintf("winget failed: %s", oneLine(string(out))))
	}
	if d.Kind == "npm" && d.NpmPkg != "" {
		log(fmt.Sprintf("npm: installing -g %s ...", d.NpmPkg))
		cmd := exec.Command("npm", "install", "-g", d.NpmPkg)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("npm install failed: %s", oneLine(string(out)))
		}
		return nil
	}
	if d.Key == "npx" {
		return fmt.Errorf("npx ships with Node.js — install Node first")
	}
	return fmt.Errorf("no automatic installer for %s on %s", d.Binary, runtime.GOOS)
}

func oneLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.Index(s, "\n"); i >= 0 {
		s = s[:i]
	}
	if len(s) > 220 {
		s = s[:220]
	}
	return s
}

// ToolkitDir returns $HOME/.vantrilex/toolkit.
func ToolkitDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".vantrilex", "toolkit"), nil
}

// ToolkitStatus reports which repos are already cloned.
// A repo counts as present when its expected dir exists with content, OR when
// any subdirectory of the toolkit root has a matching git origin remote
// (tolerates caches created under alternate directory names).
func ToolkitStatus() (root string, present map[string]bool, err error) {
	root, err = ToolkitDir()
	if err != nil {
		return "", nil, err
	}
	present = map[string]bool{}
	remotes := map[string]string{} // normalized remote url -> dirname
	if entries, err := os.ReadDir(root); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			if u := gitRemoteURL(filepath.Join(root, e.Name())); u != "" {
				remotes[normalizeGitURL(u)] = e.Name()
			}
		}
	}
	for _, r := range ToolkitRepos {
		p := filepath.Join(root, r.Dir)
		ok := dirExists(filepath.Join(p, ".git")) || dirExists(p) && hasFiles(p)
		if !ok {
			if _, found := remotes[normalizeGitURL(r.URL)]; found {
				ok = true
			}
		}
		present[r.Dir] = ok
	}
	return root, present, nil
}

// ToolkitDirName returns the on-disk directory holding a repo URL, or "".
func ToolkitDirName(url string) string {
	root, err := ToolkitDir()
	if err != nil {
		return ""
	}
	for _, r := range ToolkitRepos {
		if normalizeGitURL(r.URL) == normalizeGitURL(url) {
			p := filepath.Join(root, r.Dir)
			if dirExists(filepath.Join(p, ".git")) {
				return r.Dir
			}
		}
	}
	if entries, err := os.ReadDir(root); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			if normalizeGitURL(gitRemoteURL(filepath.Join(root, e.Name()))) == normalizeGitURL(url) {
				return e.Name()
			}
		}
	}
	return ""
}

func gitRemoteURL(dir string) string {
	if !dirExists(filepath.Join(dir, ".git")) {
		return ""