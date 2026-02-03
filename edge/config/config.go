package config

import (
	"encoding/json"
	"os"
	"strconv"
)

// these variables defined here are used up in heartbeat file
type Config struct {
	EdgeID            string `json:"edge_id"`
	Region            string `json:"region"`
	Port              int    `json:"port"`
	HeartbeatInterval int    `json:"heartbeat_interval"`
	ControlPlaneURL   string `json:"control_plane_url"`
	OriginURL         string `json:"origin_url"` // Add this
	CacheDir          string `json:"cache_dir"`  // Add this
}

func Load() *Config {
	file, err := os.Open("config.json")
	if err != nil {
		panic("failed to open config.json" + err.Error())
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		panic("failed to decode config.json: " + err.Error())
	}

	// --- DOCKER OVERRIDES START ---
	// This part allows docker-compose.yml to change settings on the fly

	if envID := os.Getenv("EDGE_ID"); envID != "" {
		cfg.EdgeID = envID
	}

	if envCP := os.Getenv("CONTROL_PLANE_URL"); envCP != "" {
		cfg.ControlPlaneURL = envCP
	}

	if envOrigin := os.Getenv("ORIGIN_URL"); envOrigin != "" {
		cfg.OriginURL = envOrigin
	}

	if envPort := os.Getenv("EDGE_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			cfg.Port = p
		}
	}
	// --- DOCKER OVERRIDES END ---

	// validations
	if cfg.EdgeID == "" {
		panic("edge_id is required")
	} else if cfg.HeartbeatInterval < 0 {
		panic("Invalid Heartbeat Interval")
	}

	return &cfg
}

// TODO: VALIDATION, ENV-VAR SUPPORT
