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