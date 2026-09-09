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
		return nil, err
	}
	out := make([]Model, 0, len(r.Data))
	for _, e := range r.Data {
		if strings.TrimSpace(e.ID) == "" {
			continue
		}
		out = append(out, entryToModel(e))
	}
	return out, nil
}

// FetchOpenRouterModels returns live models, cached per session.
func FetchOpenRouterModels(ctx context.Context) []Model {
	orCacheMu.Lock()
	if time.Since(orCachedAt) < orCacheTTL && len(orCache) > 0 {
		out := orCache
		orCacheMu.Unlock()
		return out
	}
	orCacheMu.Unlock()

	req, err := http.NewRequestWithContext(ctx, "GET", openRouterModelsURL, nil)
	if err != nil {
		return nil
	}
	resp, err := orHTTP.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil
	}
	var r openRouterResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil
	}
	models := make([]Model, 0, len(r.Data))
	for _, e := range r.Data {
		if strings.TrimSpace(e.ID) == "" {
			continue
		}
		models = append(models, entryToModel(e))
	}
	if len(models) == 0 {
		return nil
	}
	orCacheMu.Lock()
	orCache = models
	orCachedAt = time.Now()
	orCacheMu.Unlock()
	return models
}

// MergeModels overlays live entries onto the verified static matrix.
// Live pricing/context wins; static runner compat, Zen flag, and effort
// gating are preserved. Unknown live IDs are appended.
func MergeModels(static, live []Model) []Model {
	byID := map[string]int{}
	out := append([]Model(nil), static...)
	for i, m := range out {
		byID[strings.ToLower(m.ID)] = i
	}
	for _, l := range live {
		if idx, ok := byID[strings.ToLower(l.ID)]; ok {
			s := out[idx]
			if l.InputPerM != "" && l.InputPerM != "varies" {
				s.InputPerM = l.InputPerM
			}
			if l.OutputPerM != "" && l.OutputPerM != "varies" {
				s.OutputPerM = l.OutputPerM
			}
			if l.Context != "" && l.Context != "n/a" {
				s.Context = l.Context
			}
			if l.Blurb != "" {
				s.Blurb = l.Blurb
			}
			if l.Reasoning {
				s.Reasoning = true
			}
			out[idx] = s
			continue
		}
		out = append(out, l)
	}
	return out
}

// LiveModels returns merged static + cached-or-fresh live catalog.
func LiveModels(ctx context.Context) []Model {
	static := Models()
	live := FetchOpenRouterModels(ctx)
	if len(live) == 0 {
		return static
	}
	return MergeModels(static, live)
}
