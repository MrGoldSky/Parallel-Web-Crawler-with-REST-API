package fetcher

import (
	"context"
	"errors"
	"testing"
	"time"
)

type DummyFetcher struct {
    Delay time.Duration
    Err error
}

func (d *DummyFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case <-time.After(d.Delay):
    }
    if d.Err != nil {
        return nil, d.Err
    }
    return []byte("data:" + url), nil
}

func TestPoolSuccess(t *testing.T) {
    ctx := context.Background()
    df := &DummyFetcher{Delay: 10 * time.Millisecond}
    pool := NewPool(ctx, df, 3, 5)
    pool.Start()
    urls := []string{"u1", "u2", "u3", "u4"}
    for _, u := range urls {
        pool.Submit(u)
    }
    pool.Stop()

    got := map[string][]byte{}
    for res := range pool.Results() {
        if res.Err != nil {
            t.Errorf("unexpected error for %s: %v", res.URL, res.Err)
        }
        got[res.URL] = res.Body
    }
    if len(got) != len(urls) {
        t.Errorf("expected %d results, got %d", len(urls), len(got))
    }
}

func TestPoolError(t *testing.T) {
    ctx := context.Background()
    testErr := errors.New("fetch failed")
    df := &DummyFetcher{Delay: 0, Err: testErr}
    pool := NewPool(ctx, df, 2, 3)
    pool.Start()

    pool.Submit("badurl")
    pool.Stop()

    for res := range pool.Results() {
        if res.Err != testErr {
            t.Errorf("expected error %v, got %v", testErr, res.Err)
        }
    }
}