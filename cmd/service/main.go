package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	go func() {
		for {
			ctx := context.Background()
			err := client.Set(ctx, "key", "value", 0).Err()
			if err != nil {
				log.Println(err)
			} else {
				log.Println("redis get: ok")
			}
			time.Sleep(time.Second)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	log.Println("Exiting...")
}
