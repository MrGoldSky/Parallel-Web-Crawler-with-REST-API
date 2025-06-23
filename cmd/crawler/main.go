package main

import (
	"log"

	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/api"
)

func main() {
    srv := api.NewServer()
    log.Println("Starting Web Crawler API on :8080")
    if err := srv.Run(":8080"); err != nil {
        log.Fatalf("server error: %v", err)
    }
}