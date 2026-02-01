package controlplane

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func Register(url string, edgeID string, region string, port int) error {
	payload := map[string]interface{}{
		"edge_id": edgeID,
		"region":  region,
		"port":    port,
	}
	body, _ := json.Marshal(payload)

	resp, err := client.Post(url+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed: %d", resp.StatusCode)
	}
	return nil
}

func PostHeartbeat(url string, edgeID string) error {
	payload := map[string]string{"edge_id": edgeID}
	body, _ := json.Marshal(payload)

	// Use fmt.Sprintf to ensure the path is clean
	fullURL := fmt.Sprintf("%s/heartbeat", url)

	resp, err := client.Post(fullURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// IMPORTANT: Check the status code!
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat failed with status: %d", resp.StatusCode)
	}

	return nil
}
