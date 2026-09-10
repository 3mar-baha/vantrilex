# Workspace Provisioning

The Skills Gate writes only the assets checked across the five registry
stages — nothing is cloned in bulk beforehand:

- `.claude/skills/<name>/SKILL.md` — checked skills (HTTPS body or starter)
- `.claude/agents/<name>.md` — checked agent personas
- `.claude/plugins/manifest.json` — checked plugin install refs
- `.claude/hooks/` plus the event map in `.claude/settings.json`
- `opencode.json` (`"mcp"`) and `.mcp.json` (`"mcpServers"`)
- `CLAUDE.md` — appended active-component manifest

## Offline fallback behavior

Each fetch tries HTTPS twice (20s timeout, 2MB cap) and falls back to an
embedded starter template carrying the item name, source, and description.
Provisioning therefore succeeds fully offline; a network connection only
upgrades bodies to their upstream definitions. Existing files are never
overwritten.
