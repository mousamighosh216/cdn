package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	EdgeID string `json:"edge_id"`
	Region string `json:"region"`
	Port   int    `json:"port"`
}

type HeartbeatRequest struct {
	EdgeID string `json:"edge_id"`
}

func main() {
	appConfig = loadConfig()

	StartHealthMonitor()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.POST("/register", registerEdge)
	r.POST("/heartbeat", heartbeat)
	r.GET("/resolve", resolve)
	r.GET("/edges", listEdges)

	r.Run(":" + strconv.Itoa(appConfig.ServerPort))
}

func registerEdge(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	ip := c.ClientIP()

	mu.Lock()
	defer mu.Unlock()

	edges[req.EdgeID] = &Edge{
		ID:       req.EdgeID,
		IP:       ip,
		Port:     req.Port,
		Region:   req.Region,
		LastSeen: time.Now(),
		Alive:    true,
	}

	c.JSON(200, gin.H{"status": "registered"})
}

func heartbeat(c *gin.Context) {
	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	edge, ok := edges[req.EdgeID]
	if !ok {
		c.JSON(404, gin.H{"error": "edge not registered"})
		return
	}

	edge.LastSeen = time.Now()
	edge.Alive = true

	c.JSON(200, gin.H{"status": "ok"})
}

func resolve(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	for _, edge := range edges {
		if edge.Alive {
			c.JSON(200, gin.H{
				"edge_id": edge.ID,
				"ip":      edge.IP,
				"port":    edge.Port,
				"region":  edge.Region,
			})
			return
		}
	}

	c.JSON(503, gin.H{"error": "no healthy edges"})
}

func listEdges(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	c.JSON(200, edges)
}
