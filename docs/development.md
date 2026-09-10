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
