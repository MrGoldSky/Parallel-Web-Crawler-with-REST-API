package storage

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"

    _ "github.com/lib/pq"
)

type PostgresStorage struct {
    db *sql.DB
}

// NewPostgresStorage opens a connection to Postgres using the given DSN
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    // Optional: configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(0)

    // Ensure schema exists
    schema := `
CREATE TABLE IF NOT EXISTS pages (
    url TEXT PRIMARY KEY,
    title TEXT,
    data JSONB
);
CREATE INDEX IF NOT EXISTS idx_pages_title ON pages USING gin (to_tsvector('english', title));
`
    if _, err := db.Exec(schema); err != nil {
        return nil, fmt.Errorf("create schema: %w", err)
    }
    return &PostgresStorage{db: db}, nil
}

// SavePage serializes PageData to JSON and upserts into postgres
func (p *PostgresStorage) SavePage(ctx context.Context, urlStr string, data interface{}) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("marshal data: %w", err)
    }
    query := `
INSERT INTO pages (url, title, data)
VALUES ($1, $2, $3)
ON CONFLICT (url) DO UPDATE SET title = EXCLUDED.title, data = EXCLUDED.data;
`
    var title string
    if m, ok := data.(map[string]interface{}); ok {
        if t, exists := m["Title"]; exists {
            title, _ = t.(string)
        }
    }
    _, err = p.db.ExecContext(ctx, query, urlStr, title, jsonData)
    if err != nil {
        return fmt.Errorf("upsert page: %w", err)
    }
    return nil
}

// SearchPages queries pages by keyword in title or URL.
func (p *PostgresStorage) SearchPages(ctx context.Context, keyword string) ([]string, error) {
    like := "%" + keyword + "%"
    query := `
SELECT url FROM pages
WHERE title ILIKE $1 OR url ILIKE $1
LIMIT 100;
`
    rows, err := p.db.QueryContext(ctx, query, like)
    if err != nil {
        return nil, fmt.Errorf("select pages: %w", err)
    }
    defer rows.Close()

    var urls []string
    for rows.Next() {
        var u string
        if err := rows.Scan(&u); err != nil {
            return nil, err
        }
        urls = append(urls, u)
    }
    return urls, nil
}