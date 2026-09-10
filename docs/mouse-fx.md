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

## Particle budgets and reduced motion

- Hover emits 2 short-lived motes per motion event.
- Clicks detonate one radial burst; stage entries fire one perimeter wave.
- The engine hard-caps live particles (~420, ambient dropped first).
- Physics integrates on fixed 16ms ticks, independent of frame rate.
- `VANTRILEX_NO_FX=1` disables trail, burst, and perimeter effects and
  renders a static frame instead.
