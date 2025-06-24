package fetcher

import (
	"context"
	"io"
    "log"
	"net/http"
	"time"
)

// HTTPFetcher is an implementation of Fetcher using net/http
type HTTPFetcher struct {
    client *http.Client
}

// NewHTTPFetcher returns HTTPFetcher with given timeout
func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
    return &HTTPFetcher{
        client: &http.Client{Timeout: timeout},
    }
}

// Fetch downloads the content at the given URL
func (h *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        log.Printf("request creation for %q failed: %v", url, err)
        return nil, err
    }

    req.Header.Set("User-Agent", "ParallelCrawler/1.0")

    resp, err := h.client.Do(req)
    if err != nil {
        log.Printf("fetch %q failed: %v", url, err)
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}


