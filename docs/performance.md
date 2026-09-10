# Performance Notes

## Render budgets

- Registry stages render a 12-row virtual window regardless of catalog size.
- Fuzzy filtering is token-AND over precomputed lowercase haystacks (<5ms).
- The wizard fits `H-1` rows; oversized content is truncated, never scrolled.
- Footer and header are width-capped to avoid terminal wrap.

## Startup profiling notes

- Embedded registries load via `sync.Once` on first stage entry (<1ms).
- The OpenRouter fetch runs in the background; startup never blocks on it.
- Runner self-updates run concurrently in background goroutines.
- Logo rendering shells to `chafa` once per size; ASCII fallback is instant.
