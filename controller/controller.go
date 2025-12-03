package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func initCluster() {
	go monitorAgents()
}

func monitorAgents() {
	for {
		cluster.Lock()
		for _, agent := range cluster.Agents {
			agent.LastCheck = time.Now()
		}
		cluster.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func sendStartToAgent(node, id string) {
	cluster.Lock()
	agent := cluster.Agents[node]
	cluster.Unlock()

	data := map[string]string{"id": id}
	jsonData, _ := json.Marshal(data)
	http.Post(agent.URL+"/start", "application/json", bytes.NewReader(jsonData))
}

func sendStopToAgent(node, id string) {
	cluster.Lock()
	agent := cluster.Agents[node]
	cluster.Unlock()
	data := map[string]string{"id": id}
	jsonData, _ := json.Marshal(data)
	http.Post(agent.URL+"/stop", "application/json", bytes.NewReader(jsonData))
}