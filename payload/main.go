package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Payload started...")
	for {
		doWork()
		time.Sleep(3 * time.Second)
	}
}