package main

import (
	"flag"
	"fmt"
)

func main() {
	controllerURL := flag.String("controller", "http://127.0.0.1:8080", "Controller URL")
	port := flag.String("port", "9090", "Agent port")
	flag.Parse()

	fmt.Println("Starting Agent on port", *port)
	startAgent(*controllerURL, *port)
}