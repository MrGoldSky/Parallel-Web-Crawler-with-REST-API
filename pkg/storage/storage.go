package storage

import "context"

type Storage interface {
    SavePage(ctx context.Context, url string, data interface{}) error
    SearchPages(ctx context.Context, keyword string) ([]string, error)
    Clear(ctx context.Context) error
}

func NewStorage(dsn string) (Storage, error) {
    return NewPostgresStorage(dsn)
}
