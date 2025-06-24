package crawler

import (
	"context"
	"net/url"
	"path"
	"sync"

	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/fetcher"
	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/parser"
	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/storage"
)

type CrawlStats struct {
    Fetched int `json:"fetched"`
    Errors int `json:"errors"`
    Queue int `json:"in_queue"`
}

type CrawlManager struct {
    fetcher fetcher.Fetcher
    parser parser.Parser
    storage storage.Storage
    workers int
    queueSize int

    mu sync.Mutex
    visited map[string]struct{}
    pending int
    stats CrawlStats
    ctx context.Context
    cancel context.CancelFunc
    pool *fetcher.Pool
}

func NewManager(f fetcher.Fetcher, p parser.Parser, s storage.Storage, workers, queueSize int) *CrawlManager {
    ctx, cancel := context.WithCancel(context.Background())
    return &CrawlManager{
        fetcher: f,
        parser: p,
        storage: s,
        workers: workers,
        queueSize: queueSize,
        visited: make(map[string]struct{}),
        ctx: ctx,
        cancel: cancel,
    }
}

func (m *CrawlManager) Start(seeds []string, maxDepth int) {
    if m.cancel != nil {
        m.cancel()
    }
    m.ctx, m.cancel = context.WithCancel(context.Background())

    m.mu.Lock()
    m.visited = make(map[string]struct{})
    m.stats = CrawlStats{}
    m.pending = 0
    m.mu.Unlock()

    m.pool = fetcher.NewPool(m.ctx, m.fetcher, m.workers, m.queueSize)
    m.pool.Start()
    go m.runBFS(seeds, maxDepth)
}

func normalizeURL(raw string) (string, error) {
    u, err := url.Parse(raw)
    if err != nil {
        return "", err
    }
    u.Fragment = ""

    clean := path.Clean(u.Path)
    if clean == "." {
        u.Path = ""
    } else {
        u.Path = clean
    }

    return u.String(), nil
}

func (m *CrawlManager) runBFS(seeds []string, maxDepth int) {
    depths := make(map[string]int)

    // Enqueue seed URLs
    for _, raw := range seeds {
        norm, err := normalizeURL(raw)
        if err != nil {
            continue
        }
        m.mu.Lock()
        if _, seen := m.visited[norm]; !seen {
            m.visited[norm] = struct{}{}
            depths[norm] = 0
            m.pending++
            m.mu.Unlock()

            m.pool.Submit(norm)
        } else {
            m.mu.Unlock()
        }
    }

    for {
        m.mu.Lock()
        if m.pending == 0 {
            m.mu.Unlock()
            break
        }
        m.mu.Unlock()

        res, ok := <-m.pool.Results()
        if !ok {
            break
        }

        m.mu.Lock()
        m.pending--
        if res.Err != nil {
            m.stats.Errors++
            m.mu.Unlock()
            continue
        }
        m.stats.Fetched++
        m.mu.Unlock()

        data, err := m.parser.Parse(res.Body)
        if err == nil {
            _ = m.storage.SavePage(context.Background(), res.URL, data)
        }

        depth := depths[res.URL]
        if err != nil || depth >= maxDepth {
            continue
        }

        for _, link := range data.InternalLinks {
            norm, err := normalizeURL(link)
            if err != nil {
                continue
            }
            m.mu.Lock()
            if _, seen := m.visited[norm]; !seen {
                m.visited[norm] = struct{}{}
                depths[norm] = depth + 1
                m.pending++
                m.mu.Unlock()

                m.pool.Submit(norm)
            } else {
                m.mu.Unlock()
            }
        }
    }
}

func (m *CrawlManager) Stop() {
    m.cancel()
    if m.pool != nil {
        m.pool.Stop()
    }
}

func (m *CrawlManager) Stats() CrawlStats {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.stats
}

func (m *CrawlManager) StoredPages(keyword string) []string {
    pages, _ := m.storage.SearchPages(context.Background(), keyword)
    if pages == nil {
        return []string{}
    }
    return pages
}
