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
		var pe []pluginEntry
		var se []skillEntry
		var he []hookEntry
		var ae []agentEntry
		_ = json.Unmarshal(mcpJSON, &me)
		_ = json.Unmarshal(pluginsJSON, &pe)
		_ = json.Unmarshal(skillsJSON, &se)
		_ = json.Unmarshal(hooksJSON, &he)
		_ = json.Unmarshal(agentsJSON, &ae)
		for _, e := range me {
			mcps = append(mcps, RegistryItem{
				Name: e.Name, Kind: "mcp", Desc: e.Description, Tags: e.Tags,
				Source: e.SourceURL, RawURL: e.RawURL,
				Install:  e.Command + " " + strings.Join(e.Args, " "),
				Selected: defaultSelected[e.Name],
			})
		}
		for _, e := range pe {
			plugins = append(plugins, RegistryItem{
				Name: e.Name, Kind: "plugin", Desc: e.Description,
				Tags:     []string{e.Marketplace},
				Source:   e.Marketplace,
				RawURL:   e.RawURL,
				Install:  e.InstallRef,
				Selected: defaultSelected[e.Name],
			})
		}
		for _, e := range se {
			skills = append(skills, RegistryItem{
				Name: e.Name, Kind: "skill", Desc: e.Description, Tags: e.Tags,
				Source: e.Source, RawURL: e.RawURL, Install: e.Source + "@" + e.Name,
				Selected: defaultSelected[e.Name],
			})
		}
		for _, e := range he {
			hooks = append(hooks, RegistryItem{
				Name: e.Name, Kind: "hook", Desc: e.Description,
				Tags:     []string{e.Event},
				Source:   e.Event,
				RawURL:   e.RawURL,
				Install:  e.Event + ":" + e.Name,
				Selected: defaultSelected[e.Name],
			})
		}
		for _, e := range ae {
			agents = append(agents, RegistryItem{
				Name: e.Name, Kind: "agent", Desc: e.Description, Tags: e.Tags,
				Source: e.Role, RawURL: e.RawURL, Install: e.Role + "/" + e.Name,
				Selected: defaultSelected[e.Name],
			})
		}
	})
}

// MCPs returns all embedded MCP servers.
func MCPs() []RegistryItem { loadAll(); return mcps }

// Plugins returns all embedded plugins.
func Plugins() []RegistryItem { loadAll(); return plugins }

// SkillsRegistry returns all embedded skills.
func SkillsRegistry() []RegistryItem { loadAll(); return skills }

// Hooks returns all embedded hooks.
func Hooks() []RegistryItem { loadAll(); return hooks }

// Agents returns all embedded agents.
func Agents() []RegistryItem { loadAll(); return agents }

// SearchRegistry filters items by token-AND fuzzy query over name+desc+tags.
func SearchRegistry(items []RegistryItem, query string) []RegistryItem {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return items
	}
	toks := strings.Fields(q)
	out := make([]RegistryItem, 0, len(items))
	for _, it := range items {
		hay := strings.ToLower(it.Name + " " + it.Desc + " " + strings.Join(it.Tags, " ") + " " + it.Source)
		ok := true
		for _, t := range toks {
			if !strings.Contains(hay, t) {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, it)
		}
	}
	return out
}

// SelectedCount counts preselected items.
func SelectedCount(items []RegistryItem) int {
	n := 0
	for _, it := range items {
		if it.Selected {
			n++
		}
	}
	return n
}
