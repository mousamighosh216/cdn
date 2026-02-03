package main

import (
	"cdn/edge/config"
	"cdn/edge/controlplane"
	"cdn/edge/heartbeat"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	cfg := config.Load()

	hostname, err := os.Hostname()
	if err == nil && hostname != "" {
		// This will change "edge-mumbai" into "edge-mumbai-a1b2c3d4"
		cfg.EdgeID = fmt.Sprintf("%s-%s", cfg.EdgeID, hostname)
	}
	// 1. Control Plane Registration
	controlplane.Register(cfg.ControlPlaneURL, cfg.EdgeID, cfg.Region, cfg.Port)

	// 2. Background Heartbeat
	go heartbeat.Start(cfg.EdgeID, cfg.HeartbeatInterval, func(payload map[string]string) error {
		err := controlplane.PostHeartbeat(cfg.ControlPlaneURL, payload["edge_id"])
		if err != nil {
			log.Printf("HEARTBEAT ERROR: %v", err) // This will now show the 404 clearly
		}
		return err
	})

	// 3. REAL CDN LOGIC START
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Edge Case: Don't cache the root "/"
		if r.URL.Path == "/" {
			fmt.Fprintf(w, "Edge %s is active.", cfg.EdgeID)
			return
		}

		cachePath := filepath.Join("./cache", r.URL.Path)

		// A. Check if file is in local cache
		if _, err := os.Stat(cachePath); err == nil {
			log.Printf("HIT: Serving %s from cache", r.URL.Path)
			w.Header().Set("X-Cache", "HIT")
			http.ServeFile(w, r, cachePath)
			return
		}

		// B. CACHE MISS: Fetch from Origin (Port 9000)
		originURL := cfg.OriginURL + r.URL.Path
		log.Printf("MISS: Fetching %s from Origin", r.URL.Path)

		resp, err := http.Get(originURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "File not found on Origin: %s", r.URL.Path)
			return
		}
		defer resp.Body.Close()

		// C. Save to Cache and Serve at the same time
		os.MkdirAll(filepath.Dir(cachePath), 0755)
		cacheFile, _ := os.Create(cachePath)
		defer cacheFile.Close()

		w.Header().Set("X-Cache", "MISS")
		// TeeReader writes to the cacheFile while it sends data to the user
		multiWriter := io.TeeReader(resp.Body, cacheFile)
		io.Copy(w, multiWriter)
	})

	log.Printf("Edge Server started on :%d", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil))
}
