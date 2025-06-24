package main

import (
	"log"
	"os"
	"time"

	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/api"
	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/crawler"
	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/fetcher"
	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/parser"
	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/storage"
	"github.com/joho/godotenv"
)

func main() {
    _ = godotenv.Load()

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL is required")
    }
    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" {
        log.Fatal("BASE_URL is required")
    }

    // Initialize components
    stor, err := storage.NewStorage(dsn)
    if err != nil {
        log.Fatalf("storage init: %v", err)
    }
    httpF := fetcher.NewHTTPFetcher(10 * time.Second)
    htmlP, err := parser.NewParser(baseURL)
    if err != nil {
        log.Fatalf("parser init: %v", err)
    }
    mgr := crawler.NewManager(httpF, htmlP, stor, 5, 100)

    // Start server
    srv := api.NewServer(mgr, stor)
    log.Println("Starting API on :8080")
    if err := srv.Run(":8080"); err != nil {
        log.Fatalf("server error: %v", err)
    }
}