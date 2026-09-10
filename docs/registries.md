# Registry Curation

## MCP servers (1024)

Sourced from the official Anthropic MCP registry, awesome-mcp-servers, and
community servers. Each entry records a runnable command (npm or docker),
argument vector, topic tags, a one-line description, and provenance URLs.

Categories covered: databases, cloud, AI, browsing, DevOps, media,
productivity, finance, search, and communications.

## Skills, plugins, and hooks

- **Skills (320):** ECC, mattpocock, ponytail, guard-skills, Anthropic
  cybersecurity, vercel-skills, and official skills.
- **Plugins (112):** official and community marketplaces with install refs.
- **Hooks (56):** `PreToolUse`, `PostToolUse`, `PreCompact`, `Stop`,
  `SessionStart`, and `UserPromptSubmit` event mappings.

## Agent disciplines and defaults

Agents (312) span Architecture, Security, Frontend, Backend, DevOps, QA,
Data, and Mobile roles. Production defaults ship preselected:

- Agent: `Lead System Architect`
- Skills: `find-skills`, `skill-creator`, ponytail and mattpocock cores
- Plugins: `commit-commands`, `circuit-breaker-guard`, `context-primer`
- Hooks: `pre-compact-checkpoint`, `dangerous-command-guard`, `format-on-edit`
- MCP: `sequential-thinking`, `filesystem`, `fetch`, `memory`
