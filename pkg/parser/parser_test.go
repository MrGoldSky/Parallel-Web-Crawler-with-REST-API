package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
    html := `<!DOCTYPE html>
    <html><head><title>Test Page</title></head>
    <body>
        <a href="/about">About</a>
        <a href="http://external.com/page">External</a>
    </body></html>`
    parser, err := NewParser("http://example.com")
    if err != nil {
        t.Fatalf("failed to create parser: %v", err)
    }
    data, err := parser.Parse([]byte(html))
    if err != nil {
        t.Fatalf("parse error: %v", err)
    }
    if data.Title != "Test Page" {
        t.Errorf("expected title 'Test Page', got '%s'", data.Title)
    }
    if len(data.InternalLinks) != 1 || data.InternalLinks[0] != "http://example.com/about" {
        t.Errorf("internal links mismatch: %v", data.InternalLinks)
    }
    if len(data.ExternalLinks) != 1 || data.ExternalLinks[0] != "http://external.com/page" {
        t.Errorf("external links mismatch: %v", data.ExternalLinks)
    }
}