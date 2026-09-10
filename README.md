# Vantrilex Workflow Launcher
```
██╗   ██╗ █████╗ ███╗   ██╗████████╗██████╗ ██╗██╗     ███████╗██╗  ██╗
██║   ██║██╔══██╗████╗  ██║╚══██╔══╝██╔══██╗██║██║     ██╔════╝╚██╗██╔╝
██║   ██║███████║██╔██╗ ██║   ██║   ██████╔╝██║██║     █████╗   ╚███╔╝
╚██╗ ██╔╝██╔══██║██║╚██╗██║   ██║   ██╔══██╗██║██║     ██╔══╝   ██╔██╗
 ╚████╔╝ ██║  ██║██║ ╚████║   ██║   ██║  ██║██║███████╗███████╗██╔╝ ██╗
  ╚═══╝  ╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚══════╝╚══════╝╚═╝  ╚═╝
```

> Vantrilex Workflow Launcher — a native Go + Bubble Tea launcher that boots
> Claude Code, OpenCode, and OpenAI Codex sessions with live model routing,
> curated asset registries, and a cinematic deep-space terminal experience.

## Features

- **Focused providers** — OpenRouter (live catalog) and OpenCode Zen.
- **Automated capability badges** — `[REASONING]`, `[VISION]`, `[CODING]`,
  `[FREE TIER]`, `[ULTRA FAST]` derived from live metadata.
- **Enterprise registries** — 1024 MCP servers, 112 plugins, 320 skills,
  56 lifecycle hooks, 312 agent personas, embedded offline.
- **Runner self-updater** — Claude Code, OpenCode, and Codex auto-update
  in the background on every launch.
- **Zero-clone provisioning** — selected assets stream over HTTPS straight
  into the target workspace; no multi-gigabyte local clones.
- **Full mouse support** — hover stardust trails, click starbursts, wheel
  pagination, and perimeter supernova transitions at 60 FPS.
- **Keyboard-first** — every mouse action mirrors a keyboard binding.

## Quickstart

```sh
go build -o vantrilex ./cmd/vantrilex/
./vantrilex
```

Pick a runner, a model, an effort level, and a workspace. Check the agents,
skills, plugins, hooks, and MCP servers you want, then launch — the runner
takes over the terminal inside your provisioned project.

## Wizard stages

| # | Stage | Purpose |
|---|-------|---------|
| 0 | Preflight Doctor | Dependency scan plus background runner updates |
| 1 | History | Resume recent sessions or start a new project |
| 2 | Runner | Claude Code, OpenCode, or Codex CLI |
| 3 | Model Matrix | Live OpenRouter catalog with Zen models included |
| 4 | Cognitive Effort | low / medium / high / xhigh / max gating |
| 5 | Workspace | Target project directory |
| 6 | Agents | 312 personas, `Lead System Architect` preselected |
| 7 | Skills | 320 skills, `find-skills` and `skill-creator` preselected |
| 8 | Plugins | 112 plugins, commit and guard defaults preselected |
| 9 | Hooks | 56 lifecycle hooks, safety defaults preselected |
| 10 | MCP Servers | 1024 servers, reasoning and filesystem defaults |
| 11 | Skills Gate | Provision everything into the workspace |
| 12 | Launch | Hand the terminal to the runner |

## Keyboard cheat-sheet

| Keys | Action |
|------|--------|
| `Up/Down`, `j/k` | Traverse lists |
| `Space` | Toggle `[X]` / `[ ]` |
| `A` | Toggle all filtered items |
| `C` | Confirm stage |
| `/` then type | Instant fuzzy search (<5ms) |
| `Esc` | Clear search, then go back |
| `[` / `]` | Switch model category tabs |
| `1-5` | Jump to category or effort level |
| `Enter` | Select / confirm |
| `B` | Go back one stage |
| `Ctrl+C` | Quit |

## Mouse cheat-sheet