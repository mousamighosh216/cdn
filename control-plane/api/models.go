package main

import (
	"sync"
	"time"
)

type Edge struct {
	ID       string
	IP       string
	Port     int
	Region   string
	LastSeen time.Time
	Alive    bool
}

var (
	edges = make(map[string]*Edge)
	mu    sync.Mutex
)
