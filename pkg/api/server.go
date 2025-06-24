package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/MrGoldSky/Parallel-Web-Crawler-with-REST-API/pkg/crawler"
)

// Server handles HTTP requests.
type Server struct {
    router  *gin.Engine
    manager *crawler.CrawlManager
}

// NewServer builds API server with manager dependency.
func NewServer(mgr *crawler.CrawlManager) *Server {
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())
    srv := &Server{router: r, manager: mgr}
    srv.routes()
    return srv
}

func (s *Server) routes() {
    api := s.router.Group("/api")
    api.POST("/crawl/start", s.startCrawl)
    api.POST("/crawl/stop", s.stopCrawl)
    api.GET("/pages", s.listPages)
    api.GET("/stats", s.stats)
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
    pages := s.manager.StoredPages(c.Query("q"))
    c.JSON(http.StatusOK, gin.H{"pages": pages})
}

func (s *Server) stats(c *gin.Context) {
    c.JSON(http.StatusOK, s.manager.Stats())
}