package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerPort                 int `json:"server_port"`
	HeartbeatTimeoutSeconds    int `json:"heartbeat_timeout_seconds"`
	HealthCheckIntervalSeconds int `json:"health_check_interval_seconds"`
}

func LoadConfig() Config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic("cannot read config.json")
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		panic("invalid config.json format")
	}

	return cfg
}
