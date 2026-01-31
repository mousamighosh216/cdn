package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func StartServer(port string, originURL string, cacheDir string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. Construct the local path
		// Example: /thinkingskillsques.txt -> ./cache/thinkingskillsques.txt
		cachePath := filepath.Join("./cache", r.URL.Path)

		// 2. Check if file exists in cache
		if _, err := os.Stat(cachePath); err == nil {
			log.Printf("CACHE HIT: Serving %s", r.URL.Path)
			w.Header().Set("X-Cache", "HIT")
			http.ServeFile(w, r, cachePath)
			return
		}

		// 3. CACHE MISS: Fetch from Origin (Assume Origin is on Port 9000)
		originURL := "http://localhost:9000" + r.URL.Path
		log.Printf("CACHE MISS: Fetching %s from Origin", r.URL.Path)

		resp, err := http.Get(originURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "File not found on Origin: %s", r.URL.Path)
			return
		}
		defer resp.Body.Close()

		// 4. Save to Cache and Serve simultaneously
		os.MkdirAll(filepath.Dir(cachePath), 0755) // Ensure folder exists
		cacheFile, _ := os.Create(cachePath)
		defer cacheFile.Close()

		w.Header().Set("X-Cache", "MISS")
		// io.TeeReader reads from origin and writes to the file at the same time
		multiWriter := io.TeeReader(resp.Body, cacheFile)
		io.Copy(w, multiWriter)
	})

	http.ListenAndServe(":"+port, nil)
}
