package crawler

import (
	"context"
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
    m.pool = fetcher.NewPool(m.ctx, m.fetcher, m.workers, m.queueSize)
    m.pool.Start()
    go m.runBFS(seeds, maxDepth)
}

func (m *CrawlManager) runBFS(seeds []string, maxDepth int) {
    depths := make(map[string]int)
    pending := 0

    for _, u := range seeds {
        m.mu.Lock()
        if _, seen := m.visited[u]; !seen {
            m.visited[u] = struct{}{}
            depths[u] = 0
            m.stats.Queue++
            m.mu.Unlock()

            m.pool.Submit(u)
            pending++
        } else {
            m.mu.Unlock()
        }
    }

    for pending > 0 {
        select {
        case <-m.ctx.Done():
            return
        case res, ok := <-m.pool.Results():
            if !ok {
                return
            }
            pending--

            depth := depths[res.URL]
            m.mu.Lock()
            if res.Err != nil {
                m.stats.Errors++
            } else {
                m.stats.Fetched++
            }
            m.mu.Unlock()

            if res.Err == nil {
                // Parse HTML
                data, err := m.parser.Parse(res.Body)
                if err == nil {
                    _ = m.storage.SavePage(context.Background(), res.URL, data)
                    if depth < maxDepth {
                        for _, link := range data.InternalLinks {
                            m.mu.Lock()
                            if _, seen := m.visited[link]; !seen {
                                m.visited[link] = struct{}{}
                                depths[link] = depth + 1
                                m.stats.Queue++
                                m.mu.Unlock()

                                m.pool.Submit(link)
                                pending++
                            } else {
                                m.mu.Unlock()
                            }
                        }
                    }
                }
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
