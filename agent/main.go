package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

type Replica struct {
	ID   string
	PID  int
	Port int
}

var (
	replicas = make(map[string]*Replica)
	mu       sync.Mutex
)

func main() {
	go registerLoop()

	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/stop", stopHandler)

	log.Println("Agent running on :9100")
	log.Fatal(http.ListenAndServe(":9100", nil))
}

func registerLoop() {
	for {
		http.Get("http://192.168.56.102:9000/register")
		time.Sleep(3 * time.Second)
	}
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	mu.Lock()
	port := 9200 + len(replicas)
	mu.Unlock()

	cmd := exec.Command("/home/user/EP/payload/payload_bin", "-port", fmt.Sprint(port))
	cmd.Start()

	mu.Lock()
	replicas[id] = &Replica{
		ID:   id,
		PID:  cmd.Process.Pid,
		Port: port,
	}
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]int{
		"Pid":  cmd.Process.Pid,
		"Port": port,
	})
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	mu.Lock()
	rep, ok := replicas[id]
	if ok {
		exec.Command("kill", fmt.Sprint(rep.PID)).Run()
		delete(replicas, id)
	}
	mu.Unlock()

	w.Write([]byte("OK"))
}