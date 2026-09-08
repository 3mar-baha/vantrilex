// Package runner launches the chosen agent subprocess and persists sessions.
package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"vantrilex/internal/catalog"
)

// Session is one launch record in $HOME/.vantrilex/sessions.json.
type Session struct {
	At        time.Time `json:"at"`
	Runner    string    `json:"runner"`
	Model     string    `json:"model"`
	Effort    string    `json:"effort"`
	Workspace string    `json:"workspace"`
}

// sessionsFile returns the JSON path, creating the dir.
func sessionsFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".vantrilex")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "sessions.json"), nil
}

// LoadSessions reads history newest-first ( tolerant of missing file).
func LoadSessions() []Session {
	path, err := sessionsFile()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var out []Session
	if err := json.Unmarshal(data, &out); err != nil {
		return nil
	}
	// Reverse to newest-first.
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

// SaveSession appends one record (cap 50).
func SaveSession(s Session) error {
	path, err := sessionsFile()
	if err != nil {
		return err
	}
	var existing []Session
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &existing)
	}
	existing = append(existing, s)
	if len(existing) > 50 {
		existing = existing[len(existing)-50:]
	}
	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// DeleteSession removes index (newest-first indexing) from history.
func DeleteSession(newestFirstIdx int) error {
	sessions := LoadSessions()
	if newestFirstIdx < 0 || newestFirstIdx >= len(sessions) {
		return fmt.Errorf("session index out of range")
	}
	// Convert to oldest-first for storage.
	oldest := make([]Session, len(sessions))
	for i, s := range sessions {
		oldest[len(sessions)-1-i] = s
	}
	keep := oldest[:0]
	_ = keep
	var out []Session
	for i, s := range oldest {
		storageIdx := i
		targetStorage := len(oldest) - 1 - newestFirstIdx
		if storageIdx == targetStorage {
			continue
		}
		out = append(out, s)
	}
	path, err := sessionsFile()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// BuildCommand returns argv for handing over to the runner.
func BuildCommand(r catalog.Runner, modelID, effort, workspace string) (string, []string) {
	switch r.ID {
	case catalog.RunnerClaude:
		args := []string{}
		if modelID != "" {
			args = append(args, "--model", modelID)
		}
		return r.Bin, args
	case catalog.RunnerOpenCode:
		args := []string{}
		if modelID != "" {
			args = append(args, "--model", modelID)
		}
		return r.Bin, args
	case catalog.RunnerCodex:
		args := []string{}
		if modelID != "" {
			args = append(args, "--model", modelID)
		}
		if effort != "" {
			args = append(args, "--effort", effort)
		}
		return r.Bin, args
	default:
		return r.Bin, nil
	}
}

// Launch hands over the terminal: sets cwd and execs the runner attached.
// It saves the session first. This never returns on success.
func Launch(r catalog.Runner, modelID, effort, workspace string) error {
	if workspace != "" {
		if err := os.MkdirAll(workspace, 0o755); err != nil {
			return err
		}
	}
	_ = SaveSession(Session{At: time.Now(), Runner: string(r.ID), Model: modelID, Effort: effort, Workspace: workspace})
	name, args := BuildCommand(r, modelID, effort, workspace)
	bin, err := exec.LookPath(name)
	if err != nil {
		// Surface a clear English error with install hint.
		return fmt.Errorf("%s not found on PATH. Install it first (doctor stage: npm install -g ...)", name)
	}
	cmd := exec.Command(bin, args...)
	if workspace != "" {
		cmd.Dir = workspace
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

// ShellOpen opens workspace in the OS file manager (best effort).
func ShellOpen(workspace string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", workspace)
	case "darwin":
		cmd = exec.Command("open", workspace)
	default:
		cmd = exec.Command("xdg-open", workspace)
	}
	_ = cmd.Start()
}

// TruncPath shortens long paths for history display.
func TruncPath(p string, n int) string {
	p = strings.TrimSpace(p)
	if len(p) <= n {
		return p
	}
	return "..." + p[len(p)-n+3:]
}
