package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting Controller on port 8080...")
	initCluster()
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/deploy", deployHandler)
	http.HandleFunc("/scale", scaleHandler)
	http.HandleFunc("/register", registerAgent)
	log.Fatal(http.ListenAndServe(":8080", nil))
}