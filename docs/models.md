# Model Providers

## OpenRouter live ingestion

On startup the launcher fetches `https://openrouter.ai/api/v1/models` once
per session (30-minute TTL cache). Live pricing and context limits overlay
the verified static matrix; unknown live ids are appended. When offline, the
static matrix serves as the full catalog — startup never blocks on network.

## OpenCode Zen catalog

The native Zen set ships verified: `opencode/zen`, `kimi-k2.6`,
`qwen3.6-plus`, `minimax-m3`, `deepseek-v4-pro`, and `glm-5.1`. Zen entries
are pinned to the OpenCode runner and are never dropped by live merges.

## Capability badges

`ModelTags` evaluates id, description, context length, and pricing:

- `[REASONING]` — thinking/reasoning traces (R1, o-series, Opus, Zen)
- `[VISION]` — multimodal and image-capable models
- `[CODING]` — code-tuned variants (Coder, Codex, Flash, Kimi, Qwen)
- `[FREE TIER]` — zero pricing on input and output
- `[ULTRA FAST]` — low-latency or high-throughput models
