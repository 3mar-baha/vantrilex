# Windows Setup

## Terminal and fonts

1. Install Windows Terminal from the Microsoft Store.
2. Install a Nerd Font and set it as the profile font.
3. Enable 24-bit color in the profile settings.
4. Confirm mouse input is enabled for click and wheel gestures.

## PowerShell usage tips

- Run `go build -o vantrilex.exe ./cmd/vantrilex/` from `vantrilex/`.
- Launch with `.\vantrilex.exe`; paths with spaces need quotes.
- `winget install -e --id Git.Git` and the Node.js LTS id cover Doctor gaps.
- Set `$env:VANTRILEX_NO_FX = 1` to disable particle effects.
