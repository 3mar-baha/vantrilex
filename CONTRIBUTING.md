# Contributing to Vantrilex Workflow Launcher

## Workflow

1. Fork the repository and create a topic branch.
2. Keep changes small, reversible, and strictly in English.
3. Run `go vet ./...` and `go test ./...` before pushing.
4. Use Conventional Commits (`feat:`, `fix:`, `docs:`, `test:`, `chore:`).

## Dataset contributions

- Registry snapshots live under `internal/catalog/data/`.
- Regenerate with `go run ./tools/generate_registries.go`.
- Verify capacities with `go run ./tools/check`.

## UI contributions

- The wizard renders a strict vertical budget; run the layout tests.
- Mouse actions must mirror keyboard bindings.
- Respect `VANTRILEX_NO_FX=1` in every new effect.
