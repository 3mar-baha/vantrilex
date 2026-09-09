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
