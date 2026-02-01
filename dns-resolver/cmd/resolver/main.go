package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type EdgeInfo struct {
	ID    string `json:"edge_id"`
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Alive bool   `json:"alive"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. Get healthy edges from Control Plane
		resp, err := http.Get("http://localhost:8080/edges")
		if err != nil {
			http.Error(w, "Control Plane unreachable", 500)
			return
		}
		defer resp.Body.Close()

		var allEdges map[string]EdgeInfo
		if err := json.NewDecoder(resp.Body).Decode(&allEdges); err != nil {
			http.Error(w, "Failed to parse edges", 500)
			return
		}

		// Inside your Resolver's http.HandleFunc
		log.Printf("Fetched %d edges from Control Plane", len(allEdges))
		for id, edge := range allEdges {
			log.Printf("Edge %s: Alive=%v", id, edge.Alive)
		}

		// 2. Filter for ALIVE edges
		var selected *EdgeInfo
		for _, edge := range allEdges {
			if edge.Alive {
				selected = &edge
				break
			}
		}

		if selected == nil {
			http.Error(w, "No healthy edges", 503)
			return
		}

		// 3. Selection Logic (Simple: Pick the first one for now)
		// will have to update the logic
		// selected := healthyEdges[0]

		// 4. Redirect the user to the Edge
		// Example: user asks for /file.txt -> redirect to http://localhost:9001/file.txt
		// edgeAddr := fmt.Sprintf("http://%s:%d%s", "localhost", selected.Port, r.URL.Path)
		// log.Printf("Routing user to Edge: %s", selected.ID)

		// http.Redirect(w, r, edgeAddr, http.StatusTemporaryRedirect)

		target := fmt.Sprintf("http://localhost:%d%s", selected.Port, r.URL.Path)

		log.Printf("Redirecting %s to %s", r.URL.Path, target)
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	})

	log.Println("Resolver started on :7000")
	log.Fatal(http.ListenAndServe(":7000", nil))
}
