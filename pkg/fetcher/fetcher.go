package fetcher

import (
	"context"
	"sync"
)

type Fetcher interface {
    Fetch(ctx context.Context, url string) (body []byte, err error)
}

type FetchResult struct{
	URL string
	Body []byte
	Err error
}

type Pool struct{
	fetcher Fetcher
	workers int
	jobs chan string
	results chan FetchResult
	wg sync.WaitGroup
	ctx context.Context
	cancel context.CancelFunc
}

func NewPool(ctx context.Context, fetcher Fetcher, workers int, queueSize int) *Pool {
    ctx, cancel := context.WithCancel(ctx)
    return &Pool{
        fetcher: fetcher,
        workers: workers,
        jobs: make(chan string, queueSize),
        results: make(chan FetchResult, queueSize),
        ctx: ctx,
        cancel: cancel,
    }
}

// Start launches the worker goroutines
func (p *Pool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker()
    }
}

// worker is the loop each worker goroutine runs
func (p *Pool) worker() {
    defer p.wg.Done()
    for {
        select {
        case <-p.ctx.Done():
            return
        case url, ok := <-p.jobs:
            if !ok {
                return
            }
            body, err := p.fetcher.Fetch(p.ctx, url)
            p.results <- FetchResult{URL: url, Body: body, Err: err}
        }
    }
}

// Submit adds a URL to the fetch queue
func (p *Pool) Submit(url string) {
    select {
    case <-p.ctx.Done():
        return
    case p.jobs <- url:
    }
}

// Results returns a channel to receive fetch results
func (p *Pool) Results() <-chan FetchResult {
    return p.results
}

// Stop shuts down the pool and waits for workers to finish
func (p *Pool) Stop() {
    close(p.jobs)
    p.wg.Wait()
    close(p.results)
    p.cancel()
}