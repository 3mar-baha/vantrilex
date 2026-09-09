// Runner self-updater: non-blocking npm upgrades for all three runners.
package doctor

import (
	"fmt"
	"os/exec"
	"sync"
)

// RunnerUpdate describes one background npm upgrade target.
type RunnerUpdate struct {
	Key    string // claude | opencode | codex
	NpmPkg string
}

// RunnerUpdates lists the three auto-updated runners.
func RunnerUpdates() []RunnerUpdate {
	return []RunnerUpdate{
		{Key: "claude", NpmPkg: "@anthropic-ai/claude-code@latest"},
		{Key: "opencode", NpmPkg: "opencode-ai@latest"},
		{Key: "codex", NpmPkg: "@openai/codex@latest"},
	}
}

var execCommand = exec.Command

// UpdateRunner runs npm install -g <pkg>@latest (non-fatal).
func UpdateRunner(u RunnerUpdate, log func(string)) error {
	log(fmt.Sprintf("runner update: %s (%s) ...", u.Key, u.NpmPkg))
	cmd := execCommand("npm", "install", "-g", u.NpmPkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", u.Key, oneLine(string(out)))
	}
	log(fmt.Sprintf("runner update: %s up-to-date", u.Key))
	return nil
}

// UpdateAllRunners upgrades every runner concurrently; returns per-key errors.
func UpdateAllRunners(log func(string)) map[string]error {