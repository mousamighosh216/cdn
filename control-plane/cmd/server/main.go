package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cdn-project/control-plane/config"
	"github.com/cdn-project/control-plane/monitor"
	"github.com/cdn-project/control-plane/registry"
	"github.com/gin-gonic/gin"
)

var edges = make(map[string]*registry.Edge)

type RegisterRequest struct {
	EdgeID string `json:"edge_id"`
	Region string `json:"region"`
	Port   int    `json:"port"`
}

type HeartbeatRequest struct {
	EdgeID string `json:"edge_id"`
}

func main() {
	fmt.Println("Attempting to load config...")
	appConfig := config.LoadConfig()
	fmt.Printf("Config loaded! Starting server on port %d\n", appConfig.ServerPort)

	monitor.StartHealthMonitor(
		appConfig.HealthCheckIntervalSeconds,
		appConfig.HeartbeatTimeoutSeconds,
	)

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

	registry.Mu.Lock()
	defer registry.Mu.Unlock()

	edges[req.EdgeID] = &registry.Edge{
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

	registry.Mu.Lock()
	defer registry.Mu.Unlock()

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
	registry.Mu.Lock()
	defer registry.Mu.Unlock()

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
	registry.Mu.Lock()
	defer registry.Mu.Unlock()

	c.JSON(200, edges)
}
