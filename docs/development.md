# Development

## Build and test

```sh
go vet ./...
go test ./...
go build -o vantrilex ./cmd/vantrilex/
```

Benchmarks (not run by default):

```sh
go test -run=NONE -bench=. ./internal/...
```

## Dataset regeneration

```sh
go run ./tools/generate_registries.go
go run ./tools/check
```

The generator live-scrapes upstream lists once, expands curated seeds to
target capacities, and writes snapshots under `internal/catalog/data/`.
Commit the refreshed JSON when counts meet the gates.
