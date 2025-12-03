package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	cluster.Lock()
	defer cluster.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cluster)
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	type Req struct {
		Service  string `json:"service"`
		Replicas int    `json:"replicas"`
	}
	var req Req
	json.Unmarshal(body, &req)

	cluster.Lock()
	defer cluster.Unlock()
	for i := 0; i < req.Replicas; i++ {
		id := fmt.Sprintf("%s-%d", req.Service, rand.Intn(10000))
		node := pickNode()
		if node == "" {
			fmt.Println("No alive agents")
			continue
		}
		replica := &Replica{ID: id, Node: node, StartTime: time.Now(), Status: "running"}
		cluster.Services[req.Service] = append(cluster.Services[req.Service], replica)
		go sendStartToAgent(node, id)
	}

	w.Write([]byte("Deploy request accepted"))
}

func scaleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Scale endpoint - placeholder"))
}

func registerAgent(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var agent AgentStatus
	json.Unmarshal(body, &agent)
	agent.Alive = true
	agent.LastCheck = time.Now()

	cluster.Lock()
	cluster.Agents[agent.Node] = &agent
	cluster.Unlock()
	fmt.Println("Registered agent:", agent.Node)
	w.Write([]byte("ok"))
}

func pickNode() string {
	for node, agent := range cluster.Agents {
		if agent.Alive {
			return node
		}
	}
	return ""
}