// Package catalog registry: embedded offline datasets for agents, skills,
// plugins, hooks, and MCP servers with instant fuzzy search.
package catalog

import (
	_ "embed"
	"encoding/json"
	"strings"
	"sync"
)

//go:embed data/mcp_registry.json
var mcpJSON []byte

//go:embed data/plugins_registry.json
var pluginsJSON []byte

//go:embed data/skills_registry.json
var skillsJSON []byte

//go:embed data/hooks_registry.json
var hooksJSON []byte

//go:embed data/agents_registry.json
var agentsJSON []byte

// RegistryItem is the unified selectable asset shape.
type RegistryItem struct {
	Name     string   `json:"name"`
	Kind     string   `json:"kind,omitempty"`
	Desc     string   `json:"description"`
	Tags     []string `json:"tags"`
	Source   string   `json:"source"`
	RawURL   string   `json:"rawURL"`
	Install  string   `json:"install"`
	Selected bool     `json:"-"`
}

type mcpEntry struct {
	Name        string   `json:"name"`
	Command     string   `json:"command"`
	Args        []string `json:"args"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	SourceURL   string   `json:"sourceURL"`
	RawURL      string   `json:"rawURL"`
}

type pluginEntry struct {
	Name        string `json:"name"`
	Marketplace string `json:"marketplace"`
	Description string `json:"description"`
	InstallRef  string `json:"installRef"`
	RawURL      string `json:"rawURL"`
}

type skillEntry struct {
	Name        string   `json:"name"`
	Source      string   `json:"source"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	RawURL      string   `json:"rawURL"`
}

type hookEntry struct {
	Name        string `json:"name"`
	Event       string `json:"event"`
	Description string `json:"description"`
	RawURL      string `json:"rawURL"`
}

type agentEntry struct {
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	RawURL      string   `json:"rawURL"`
}

var (
	loadOnce sync.Once
	mcps     []RegistryItem
	plugins  []RegistryItem
	skills   []RegistryItem
	hooks    []RegistryItem
	agents   []RegistryItem
)

var defaultSelected = map[string]bool{
	// Agents
	"Lead System Architect": true,
	// Skills
	"find-skills": true, "skill-creator": true, "ponytail-core": true, "mattpocock-typescript": true,
	// Plugins
	"commit-commands": true, "circuit-breaker-guard": true, "context-primer": true,
	// Hooks
	"pre-compact-checkpoint": true, "dangerous-command-guard": true, "format-on-edit": true,
	// MCP
	"sequential-thinking": true, "filesystem": true, "fetch": true, "memory": true,
}

func loadAll() {
	loadOnce.Do(func() {
		var me []mcpEntry