package main

import (
	"log"
	"time"
)

func main() {
	for i := 1; i <= 10; i++ {
		log.Printf("Doing work (%d) ...", i)
		time.Sleep(time.Second)
	}
	log.Println("Exiting...")
}
