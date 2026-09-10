# System Architecture

Vantrilex Workflow Launcher is a native Go binary built on Bubble Tea
(model-update-view) with Lipgloss styling. There is no daemon, no database,
and no network dependency at runtime beyond on-demand asset fetching.

```
┌──────────┐   tea.Msg    ┌──────────┐   exec   ┌──────────────┐
│ Bubble   │ ──────────▶  │ Update   │ ───────▶ │ Runner CLI   │
│ Tea loop │ ◀──────────  │ handlers │          │ (handover)   │
└──────────┘   60 FPS     └──────────┘          └──────────────┘
     │ tick                │ stages             workspace
     ▼                     ▼                    provisioned
┌──────────┐         ┌──────────┐          ┌──────────────┐
│ Particles│         │ Catalog  │          │ .claude/ +   │
│ engine   │         │ registry │          │ opencode.json│
└──────────┘         └──────────┘          └──────────────┘
```

State flows one way: messages update the model, the model renders the view,
and the launcher hands the terminal to the chosen runner on confirm.

## Wizard stage flow

```
Doctor → History → Runner → Model → Effort → Workspace
  → Agents → Skills → Plugins → Hooks → MCP → Gate → Launch
```

- Forward motion fires a perimeter supernova through the particle engine.
- Every stage fits a strict vertical budget (`H-1` rows) so nothing scrolls.
- Registry stages share one virtualized list component with a 12-row window.
- The Skills Gate is the only mutating stage: it provisions the workspace.
