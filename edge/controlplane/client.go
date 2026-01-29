package controlplane

import "net/http"

func PostHeartbeat(url string, payload map[string]string) error {
	_, err := http.Post(url+"/heartbeat", "application/json", nil)
	return err
}
