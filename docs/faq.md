# Troubleshooting FAQ

## The runner is "not found on PATH" at launch

Install it from the Doctor stage (`A` auto-installs) or manually:

```sh
npm install -g @anthropic-ai/claude-code@latest
npm install -g opencode-ai@latest
npm install -g @openai/codex@latest
```

## Live models never load

The OpenRouter fetch needs outbound HTTPS. Offline, the verified static
matrix (including the full Zen set) is used automatically.

## Provisioning wrote starter templates instead of upstream bodies

The HTTPS fetch failed or was offline. Re-run provisioning with network
access after deleting the starter files you want upgraded.

## Windows terminal setup

- Prefer Windows Terminal with a Nerd Font for box glyphs (`▸ ⟦⟧ ✦`).
- Enable mouse input in the terminal profile for click and wheel support.
- If colors look flat, set the profile to 24-bit ("true color") mode.

## PowerShell tips

- Launch with `.antrilex.exe` from the project root.
- Quote workspace paths containing spaces.
- `Ctrl+C` always exits the wizard safely without side effects.
