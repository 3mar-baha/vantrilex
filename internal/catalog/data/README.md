# Embedded Dataset Snapshots

Pre-indexed offline JSON registries embedded into the binary via
`//go:embed`. Regenerate with `go run ./tools/generate_registries.go`
and validate with `go run ./tools/check`.

| File | Entries | Gate |
|------|---------|------|
| `mcp_registry.json` | 1024 | >=1000 |
| `plugins_registry.json` | 112 | >=100 |
| `skills_registry.json` | 320 | >=300 |
| `hooks_registry.json` | 56 | >=50 |
| `agents_registry.json` | 312 | >=300 |
