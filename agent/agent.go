package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sync"
)

var mu sync.Mutex
var replicas = make(map[string]*exec.Cmd)

func startAgent(controllerURL, port string) {
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/stop", handleStop)
	http.HandleFunc("/status", handleStatus)

	fmt.Println("Agent running on port", port)
	go registerToController(controllerURL, port)

	http.ListenAndServe(":"+port, nil)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		ID string `json:"id"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var req Req
	json.Unmarshal(body, &req)
	startPayload(req.ID)
	w.Write([]byte("Started payload " + req.ID))
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		ID string `json:"id"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var req Req
	json.Unmarshal(body, &req)
	stopPayload(req.ID)
	w.Write([]byte("Stopped payload " + req.ID))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	status := make(map[string]string)
	for id := range replicas {
		status[id] = "running"
	}
	json.NewEncoder(w).Encode(status)
}

func startPayload(id string) {
	mu.Lock()
	defer mu.Unlock()
	cmd := exec.Command("go", "run", "../payload/main.go")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to start payload:", err)
		return
	}
	replicas[id] = cmd
	fmt.Println("Payload started with ID:", id)
}

func stopPayload(id string) {
	mu.Lock()
	defer mu.Unlock()
	if cmd, ok := replicas[id]; ok {
		cmd.Process.Kill()
		delete(replicas, id)
		fmt.Println("Payload stopped with ID:", id)
	}
}

func registerToController(controllerURL, port string) {
	for {
		data := map[string]string{"node": port, "url": "http://127.0.0.1:" + port}
		jsonData, _ := json.Marshal(data)
		http.Post(controllerURL+"/register", "application/json", bytes.NewReader(jsonData))
	}
}