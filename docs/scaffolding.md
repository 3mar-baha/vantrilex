# Workspace Provisioning

The Skills Gate writes only the assets checked across the five registry
stages — nothing is cloned in bulk beforehand:

- `.claude/skills/<name>/SKILL.md` — checked skills (HTTPS body or starter)
- `.claude/agents/<name>.md` — checked agent personas
- `.claude/plugins/manifest.json` — checked plugin install refs
- `.claude/hooks/` plus the event map in `.claude/settings.json`
- `opencode.json` (`"mcp"`) and `.mcp.json` (`"mcpServers"`)
- `CLAUDE.md` — appended active-component manifest
