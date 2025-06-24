package api

import (
	"context"
	"net/http"

	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/crawler"
	"github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/storage"
	"github.com/gin-gonic/gin"
)

// Server handles HTTP requests.
type Server struct {
    router  *gin.Engine
    manager *crawler.CrawlManager
    storage storage.Storage
}

// NewServer builds API server with manager dependency.
func NewServer(mgr *crawler.CrawlManager, storage storage.Storage) *Server {
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())
    srv := &Server{router: r, manager: mgr, storage: storage}
    srv.routes()
    return srv
}

func (s *Server) routes() {
    api := s.router.Group("/api")
    api.POST("/crawl/start", s.startCrawl)
    api.POST("/crawl/stop", s.stopCrawl)
    api.GET("/pages", s.listPages)
    api.GET("/stats", s.stats)
    api.DELETE("/pages", s.clearPages)
}

func (s *Server) Run(addr string) error {
    return s.router.Run(addr)
}

type crawlRequest struct {
    Seeds []string `json:"seeds" binding:"required"`
    MaxDepth int `json:"max_depth" binding:"required"`
}

func (s *Server) startCrawl(c *gin.Context) {
    var req crawlRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    s.manager.Start(req.Seeds, req.MaxDepth)
    c.JSON(http.StatusOK, gin.H{"status": "started"})
}

func (s *Server) stopCrawl(c *gin.Context) {
    s.manager.Stop()
    c.JSON(http.StatusOK, s.manager.Stats())
}

func (s *Server) listPages(c *gin.Context) {
    q := c.Query("q")
    pages, err := s.storage.SearchPages(context.Background(), q)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if pages == nil {
        pages = []string{}
    }
    c.JSON(http.StatusOK, gin.H{"pages": pages})
}

func (s *Server) stats(c *gin.Context) {
    c.JSON(http.StatusOK, s.manager.Stats())
}

func (s *Server) clearPages(c *gin.Context) {
    if err := s.storage.Clear(context.Background()); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "cleared"})
}
