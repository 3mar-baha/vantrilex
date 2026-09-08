// Live OpenRouter model ingestion with per-session cache and merge.
package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const openRouterModelsURL = "https://openrouter.ai/api/v1/models"

type openRouterPricing struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}

type openRouterEntry struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	ContextLength int               `json:"context_length"`
	Pricing       openRouterPricing `json:"pricing"`
}

type openRouterResp struct {
	Data []openRouterEntry `json:"data"`
}

var (
	orCacheMu  sync.Mutex
	orCache    []Model
	orCachedAt time.Time
	orCacheTTL = 30 * time.Minute
	orHTTP     = &http.Client{Timeout: 15 * time.Second}
)

func formatPrice(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return "varies"
	}
	var f float64
	if _, err := fmt.Sscanf(p, "%f", &f); err != nil {
		return p
	}