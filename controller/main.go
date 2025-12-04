package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AgentInfo struct {
	Address         string
	LastSeen        time.Time
	Active          bool
	RunningReplicas int
}

type Replica struct {
	ID    string
	Agent string
	PID   int
	Start time.Time
	Port  int
}

type ClusterState struct {
	Agents   map[string]*AgentInfo
	Replicas map[string]*Replica
	Desired  int
}

var (
	state = ClusterState{
		Agents:   make(map[string]*AgentInfo),
		Replicas: make(map[string]*Replica),
		Desired:  2,
	}
	mu sync.Mutex
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/scale", handleScale)
	http.HandleFunc("/proxy", handleProxy)

	go monitorAgents()
	go reconcile()

	log.Println("Controller running on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr[:len(r.RemoteAddr)-6]

	mu.Lock()
	defer mu.Unlock()

	if _, ok := state.Agents[ip]; !ok {
		state.Agents[ip] = &AgentInfo{
			Address:         ip,
			LastSeen:        time.Now(),
			Active:          true,
			RunningReplicas: 0,
		}
		log.Println("New agent:", ip)
	}

	state.Agents[ip].LastSeen = time.Now()
	state.Agents[ip].Active = true
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	json.NewEncoder(w).Encode(state)
}

func handleScale(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Replicas int `json:"replicas"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	mu.Lock()
	state.Desired = req.Replicas
	mu.Unlock()

	w.Write([]byte("OK"))
}

func reconcile() {
	for {
		time.Sleep(2 * time.Second)

		mu.Lock()
		cur := len(state.Replicas)
		desired := state.Desired

		if cur < desired {
			need := desired - cur
			for i := 0; i < need; i++ {
				startReplica()
			}
		}

		if cur > desired {
			extra := cur - desired
			for i := 0; i < extra; i++ {
				stopReplica()
			}
		}

		mu.Unlock()
	}
}

func pickAgent() string {
	agents := make([]*AgentInfo, 0)
	for _, a := range state.Agents {
		if a.Active {
			agents = append(agents, a)
		}
	}
	if len(agents) == 0 {
		return ""
	}

	sort.Slice(agents, func(i, j int) bool {
		return agents[i].RunningReplicas < agents[j].RunningReplicas
	})

	return agents[0].Address
}

func startReplica() {
	agent := pickAgent()
	if agent == "" {
		log.Println("NO ACTIVE AGENTS")
		return
	}

	id := uuid.New().String()

	resp, err := http.Get(fmt.Sprintf("http://%s:9100/start?id=%s", agent, id))
	if err != nil {
		log.Println("Start error:", err)
		return
	}

	var data struct {
		Pid  int
		Port int
	}
	json.NewDecoder(resp.Body).Decode(&data)

	state.Replicas[id] = &Replica{
		ID:    id,
		Agent: agent,
		PID:   data.Pid,
		Port:  data.Port,
		Start: time.Now(),
	}

	state.Agents[agent].RunningReplicas++
	log.Println("Replica started:", id)
}

func stopReplica() {
	for id, rep := range state.Replicas {
		http.Get(fmt.Sprintf("http://%s:9100/stop?id=%s", rep.Agent, rep.ID))
		state.Agents[rep.Agent].RunningReplicas--
		delete(state.Replicas, id)
		return
	}
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	replicas := make([]*Replica, 0)
	for _, rep := range state.Replicas {
		replicas = append(replicas, rep)
	}
	mu.Unlock()

	if len(replicas) == 0 {
		w.WriteHeader(500)
		w.Write([]byte("NO REPLICAS"))
		return
	}

	rep := replicas[rand.Intn(len(replicas))]
	url := fmt.Sprintf("http://%s:%d/work", rep.Agent, rep.Port)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Replica error"))
		return
	}

	data := make([]byte, 8)
	n, _ := resp.Body.Read(data)
	w.Write(data[:n])
}

func monitorAgents() {
	for {
		time.Sleep(3 * time.Second)

		mu.Lock()
		now := time.Now()
		for _, a := range state.Agents {
			if now.Sub(a.LastSeen) > 10*time.Second {
				a.Active = false
			}
		}
		mu.Unlock()
	}
}