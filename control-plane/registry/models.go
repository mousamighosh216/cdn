package registry

import (
	"sync"
	"time"
)

type Edge struct {
	ID       string    `json:"id"`
	IP       string    `json:"ip"`
	Port     int       `json:"port"`
	Region   string    `json:"region"`
	LastSeen time.Time `json:"last_seen"`
	Alive    bool      `json:"alive"`
}

var (
	Data = make(map[string]*Edge)
	Mu   sync.Mutex // Capitalized so other packages can lock/unlock
)
