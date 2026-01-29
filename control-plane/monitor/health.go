package monitor

import (
	"time"

	"github.com/cdn-project/control-plane/registry"
)

func StartHealthMonitor(intervalSeconds int, timeoutSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)

	go func() {
		for range ticker.C {
			registry.Mu.Lock()
			now := time.Now()

			for _, edge := range registry.Data {
				// If the difference between 'now' and 'LastSeen' is too big...
				if now.Sub(edge.LastSeen) > time.Duration(timeoutSeconds)*time.Second {
					edge.Alive = false
				}
			}
			registry.Mu.Unlock()
		}
	}()
}
