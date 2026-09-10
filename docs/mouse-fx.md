# Mouse Interaction Model

Mouse support is additive: every gesture mirrors a keyboard binding, and the
wizard remains fully operable without a mouse.

| Gesture      | Effect |
|--------------|--------|
| Hover        | Cyan/indigo micro-particle trail (`· ⋆ ∘`) |
| Left click   | 12–16 starburst (`✦ ✧ ⚡`) plus instant select/toggle |
| Wheel        | Paginate the active virtual list |
| Right click  | Cycle model category tabs |
| Click search | Focus the workspace input |

Bubble Tea runs with `tea.WithMouseCellMotion()` for per-cell tracking.
