package heartbeat

import "time"

type Sender func(payload map[string]string) error

func Start(
	edgeID string,
	intervalSeconds int,
	send Sender,
) {
	// creates a clk that sends a signal on a channel at regular intervals
	ticker := time.NewTicker(time.Duration((intervalSeconds)) * time.Second)

	// goroutine
	go func() {
		// pauses until ticker sends signal
		for range ticker.C {
			send(map[string]string{
				"edge_id": edgeID,
			})
		}
	}()
}
