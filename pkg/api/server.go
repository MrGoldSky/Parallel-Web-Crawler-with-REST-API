package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Server struct {
    router *gin.Engine
}

func NewServer() *Server {
    r := gin.Default()
    s := &Server{router: r}
    s.registerRoutes()
    return s
}

func (s *Server) registerRoutes() {
    api := s.router.Group("/api")
    api.POST("/crawl/start", s.startCrawl)
    api.POST("/crawl/stop", s.stopCrawl)
    api.GET("/pages", s.listPages)
    api.GET("/stats", s.stats)
}

func (s *Server) Run(addr string) error {
    return s.router.Run(addr)
}

func (s *Server) startCrawl(c *gin.Context) {
    // TODO: implement
    c.JSON(http.StatusOK, gin.H{"status": "started"})
}

func (s *Server) stopCrawl(c *gin.Context) {
    // TODO: implement
    c.JSON(http.StatusOK, gin.H{"status": "stopped"})
}

func (s *Server) listPages(c *gin.Context) {
    // TODO: implement
    c.JSON(http.StatusOK, gin.H{"pages": []string{}})
}

func (s *Server) stats(c *gin.Context) {
    // TODO: implement
    c.JSON(http.StatusOK, gin.H{"stats": gin.H{}})
}
