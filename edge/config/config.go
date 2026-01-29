package config

import (
	"encoding/json"
	"os"
)

// these variables defined here are used up in heartbeat file
type Config struct {
	EdgeID            string `json:"edge_id"`
	Region            string `json:"region"`
	Port              int    `json:"port"`
	HeartbeatInterval int    `json:"heartbeat_interval"`
	ControlPlaneURL   string `json:"control_plane_url"`
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

	// validations
	if cfg.EdgeID == "" {
		panic("edge_id is required")
	} else if cfg.HeartbeatInterval < 0 {
		panic("Invalid Heartbeat Interval")
	}

	return &cfg
}

// TODO: VALIDATION, ENV-VAR SUPPORT
