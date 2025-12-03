package main

import (
	"sync"
	"time"
)

type ClusterState struct {
	sync.Mutex
	Services map[string][]*Replica
	Agents   map[string]*AgentStatus
}

type Replica struct {
	ID        string
	Node      string
	StartTime time.Time
	Status    string
}

type AgentStatus struct {
	Node      string
	URL       string
	Alive     bool
	LastCheck time.Time
}

var cluster = &ClusterState{
	Services: make(map[string][]*Replica),
	Agents:   make(map[string]*AgentStatus),
}