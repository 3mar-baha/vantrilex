// On-demand asset fetcher: zero bulk clones. Selected registry items are
// fetched over HTTPS only when needed and written straight into the target
// project workspace. Offline-safe: starter bodies are embedded so Apply
// succeeds without network; HTTPS upgrades content best-effort.
package scaffold

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vantrilex/internal/catalog"
)

var fetchHTTP = &http.Client{Timeout: 20 * time.Second}

// FetchBody downloads a raw definition over HTTPS (single retry).
func FetchBody(ctx context.Context, url string) ([]byte, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil, fmt.Errorf("empty url")
	}
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := fetchHTTP.Do(req)
		if err != nil {
			last = err
			continue
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		resp.Body.Close()
		if err != nil {
			last = err
			continue
		}
		if resp.StatusCode != 200 {
			last = fmt.Errorf("http %d for %s", resp.StatusCode, url)
			continue
		}
		if len(body) == 0 {
			last = fmt.Errorf("empty body for %s", url)
			continue
		}
		return body, nil
	}
	return nil, last
}

func writeOnceFile(abs, content string, created *[]string, rel string) error {
	if _, err := os.Stat(abs); err == nil {
		return nil
	}