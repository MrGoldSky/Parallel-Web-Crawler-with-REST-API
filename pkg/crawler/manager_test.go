package crawler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/parser"
)

type stubFetcher struct {
    responses [][]byte
    idx       int
}

func (s *stubFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
    if s.idx >= len(s.responses) {
        return nil, errors.New("no more responses")
    }
    resp := s.responses[s.idx]
    s.idx++
    return resp, nil
}

type stubParser struct{}

func (p *stubParser) Parse(html []byte) (parser.PageData, error) {
    return parser.PageData{
        Title:         string(html),
        InternalLinks: []string{"http://example.com/next"},
        ExternalLinks: nil,
    }, nil
}

type stubStorage struct {
    saved map[string]interface{}
	pages []string
}

func newStubStorage() *stubStorage {
    return &stubStorage{saved: make(map[string]interface{})}
}

func (s *stubStorage) SavePage(ctx context.Context, url string, data interface{}) error {
    s.saved[url] = data
    return nil
}

func (s *stubStorage) SearchPages(ctx context.Context, keyword string) ([]string, error) {
    var list []string
    for url := range s.saved {
        if keyword == "" || contains(url, keyword) {
            list = append(list, url)
        }
    }
    return list, nil
}

func contains(str, substr string) bool {
    return len(substr) == 0 || (len(str) >= len(substr) && str[:len(substr)] == substr)
}

func (s *stubStorage) Clear(ctx context.Context) error {
    s.pages = make([]string, 0)
    return nil
}

func TestCrawlManagerBFS(t *testing.T) {
    fetcher := &stubFetcher{responses: [][]byte{[]byte("page1"), []byte("page2")}}
    parser := &stubParser{}
    stor := newStubStorage()
    mgr := NewManager(fetcher, parser, stor, 2, 10)

    mgr.Start([]string{"http://example.com"}, 1)
    time.Sleep(100 * time.Millisecond)
    mgr.Stop()

    stats := mgr.Stats()
    if stats.Fetched != 2 {
        t.Errorf("expected fetched=2, got %d", stats.Fetched)
    }
    pages, _ := stor.SearchPages(context.Background(), "")
    if len(pages) != 2 {
        t.Errorf("expected 2 saved pages, got %d", len(pages))
    }
}