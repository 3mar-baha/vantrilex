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
	return fmt.Sprintf("$%.2f", f*1_000_000)
}

func formatContext(n int) string {
	if n <= 0 {
		return "n/a"
	}
	if n >= 1000 && n%1000 == 0 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return fmt.Sprintf("%d", n)
}

func entryToModel(e openRouterEntry) Model {
	short := e.Name
	if short == "" {
		short = e.ID
	}
	if i := strings.Index(short, " ("); i > 0 {
		short = short[:i]
	}
	if len(short) > 42 {
		short = short[:42]
	}
	m := Model{
		ID:         e.ID,
		Short:      short,
		Category:   CatAll,
		InputPerM:  formatPrice(e.Pricing.Prompt),
		OutputPerM: formatPrice(e.Pricing.Completion),
		Context:    formatContext(e.ContextLength),
		Latency:    "MED",
		Blurb:      truncate(e.Description, 90),
	}
	m.Reasoning = hasReasoningSignal(e.ID + " " + e.Description)
	return m
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

// ParseOpenRouter payload helper (testable, no network).
func ParseOpenRouter(data []byte) ([]Model, error) {
	var r openRouterResp
	if err := json.Unmarshal(data, &r); err != nil {