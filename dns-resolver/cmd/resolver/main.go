package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type EdgeInfo struct {
	ID    string `json:"edge_id"`
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Alive bool   `json:"alive"`
}

// Global counter that lives outside the handler
// Using uint64 for atomic operations
var requestCounter uint64

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var allEdges map[string]EdgeInfo
		var healthyEdges []EdgeInfo

		// 1. Get edges from Control Plane
		resp, err := http.Get("http://localhost:8080/edges")
		if err != nil {
			http.Error(w, "Control Plane unreachable", 500)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&allEdges); err != nil {
			http.Error(w, "Failed to parse edges", 500)
			return
		}

		// 2. Filter for ALIVE edges
		for _, edge := range allEdges {
			if edge.Alive {
				healthyEdges = append(healthyEdges, edge)
			}
		}

		numHealthy := len(healthyEdges)
		if numHealthy == 0 {
			log.Println("Wait! No edges registered yet.")
			http.Error(w, "CDN Warming Up - No Edges Available", 503)
			return
		}

		// 🎡 ROUND ROBIN - THE ATOMIC WAY
		// This safely increments and gets the new number across all users
		count := atomic.AddUint64(&requestCounter, 1)

		// Use modulo to pick the edge
		index := int((count - 1) % uint64(numHealthy))
		selected := healthyEdges[index]

		// Construct target (Using localhost for Docker port mapping simulation)
		target := fmt.Sprintf("http://localhost:%d%s", selected.Port, r.URL.Path)

		log.Printf("[Request %d] Redirecting to Edge: %s (Port: %d)", count, selected.ID, selected.Port)
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	})

	log.Println("Resolver started on :7000")
	log.Fatal(http.ListenAndServe(":7000", nil))
}
