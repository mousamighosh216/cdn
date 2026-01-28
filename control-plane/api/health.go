package main

import "time"

func StartHealthMonitor() {
	go func() {
		for {
			time.Sleep(time.Duration(appConfig.HealthCheckIntervalSeconds) * time.Second)
			checkEdgesHealth()
		}
	}()
}

func checkEdgesHealth() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	for _, edge := range edges {
		if now.Sub(edge.LastSeen) >
			time.Duration(appConfig.HeartbeatTimeoutSeconds)*time.Second {
			edge.Alive = false
		}
	}
}
