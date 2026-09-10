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
