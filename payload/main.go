package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var port int
	flag.IntVar(&port, "port", 8080, "")
	flag.Parse()

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(200)+50) * time.Millisecond)
		w.Write([]byte("ok"))
	})

	log.Println("Payload running on", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}