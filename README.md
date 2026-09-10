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

| Gesture | Action |
|---------|--------|
| Hover | Delicate cyan/indigo stardust trail |
| Left click | Starburst plus instant select / toggle |
| Wheel up/down | Paginate virtualized lists |
| Right click | Cycle model category tabs |
| Click search bar | Focus workspace input |

Set `VANTRILEX_NO_FX=1` to disable particle effects (static frame).

## Model providers and tagging

- **OpenRouter** — the complete live catalog at
  `https://openrouter.ai/api/v1/models`, cached per session with the
  verified static matrix as an offline fallback.
- **OpenCode Zen** — the native Zen set (`opencode/zen`, `kimi-k2.6`,
  `qwen3.6-plus`, `minimax-m3`, `deepseek-v4-pro`, `glm-5.1`).
- Badges evaluate model id, description, context length, and pricing, and
  every row shows live pricing per 1M input/output tokens plus context.

## Embedded catalogs

Pre-indexed offline JSON datasets under `internal/catalog/data/`:

| Dataset | Count | Sources |
|---------|-------|---------|
| MCP servers | 1024 | Official registry, awesome lists, community |
| Plugins | 112 | Official and community marketplaces |
| Skills | 320 | ECC, mattpocock, ponytail, guard-skills, vercel |
| Hooks | 56 | `PreToolUse`, `PostToolUse`, `PreCompact`, `Stop` |
| Agents | 312 | Architecture, Security, Frontend, Backend, DevOps, QA |

Regenerate with `go run ./tools/generate_registries.go` (live-scrapes once,
then embeds the snapshot).

## Workspace outputs

Provisioning writes only the assets you checked:

- `.claude/skills/<name>/SKILL.md`
- `.claude/agents/<name>.md`
- `.claude/plugins/manifest.json`
- `.claude/hooks/` plus event map in `.claude/settings.json`
- `opencode.json` (`"mcp"`) and `.mcp.json` (`"mcpServers"`)
- `CLAUDE.md` with the exhaustive active-component manifest

## Development

```sh
go vet ./...
go test ./...
go build -o vantrilex ./cmd/vantrilex/
```

The particle engine ticks at a fixed 16ms step, lists render a 12-row
virtual window, and all UI copy is strictly English.

## Project structure

```
cmd/vantrilex/        entrypoint with mouse cell-motion
internal/catalog/     runners, live models, tags, embedded registries
internal/doctor/      preflight checks plus runner self-updater
internal/particles/   60 FPS stardust, bursts, perimeter shockwaves
internal/runner/      session history and subprocess handover
internal/scaffold/    on-demand HTTPS fetcher and workspace writer
internal/ui/          Bubble Tea wizard, virtual lists, mouse system
tools/                registry dataset generator
```

## License

MIT — see `LICENSE` (to be added with the next release train).

## Documentation index

- `docs/architecture.md` — system overview and wizard flow.
- `docs/models.md` — providers, Zen catalog, badges, effort gating.
- `docs/registries.md` — embedded catalog curation and defaults.
- `docs/scaffolding.md` — workspace provisioning and offline fallback.
- `docs/mouse-fx.md` — mouse model and particle budgets.
- `docs/performance.md` — render budgets and profiling notes.
- `docs/windows.md` — terminal setup on Windows.
- `docs/development.md` — build, test, and dataset commands.
- `docs/roadmap.md` — v1.1 milestones and release train.
- `docs/faq.md` — troubleshooting answers.

## Status

![ci](https://github.com/3mar-baha/vantrilex/actions/workflows/ci.yml/badge.svg)
![release](https://img.shields.io/github/v/release/3mar-baha/vantrilex)
![go](https://img.shields.io/badge/go-1.24-00ADD8)
![license](https://img.shields.io/badge/license-MIT-green)
