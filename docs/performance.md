# Performance Notes

## Render budgets

- Registry stages render a 12-row virtual window regardless of catalog size.
- Fuzzy filtering is token-AND over precomputed lowercase haystacks (<5ms).
- The wizard fits `H-1` rows; oversized content is truncated, never scrolled.
- Footer and header are width-capped to avoid terminal wrap.
