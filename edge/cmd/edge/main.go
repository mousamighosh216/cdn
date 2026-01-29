package main

import (
	"cdn/edge/config"
	"cdn/edge/controlplane"
	"cdn/edge/heartbeat"
)

func main() {
	cfg := config.Load()

	heartbeat.Start(
		cfg.EdgeID,
		cfg.HeartbeatInterval,
		func(payload map[string]string) error {
			return controlplane.PostHeartbeat(cfg.ControlPlaneURL, payload)
		},
	)

	// ensures even if background is progressing main func doesnt shuts down
	select {}
}
