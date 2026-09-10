# Model Providers

## OpenRouter live ingestion

On startup the launcher fetches `https://openrouter.ai/api/v1/models` once
per session (30-minute TTL cache). Live pricing and context limits overlay
the verified static matrix; unknown live ids are appended. When offline, the
static matrix serves as the full catalog — startup never blocks on network.
