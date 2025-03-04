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

// stressTestRedis отправит n кол-во set запросов в Редис.
func stressTestRedis(client *redis.Client, n int) error {
	for i := 0; i < n; i++ {
		err := client.Set(context.Background(), "key", "value", 0).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Println("[ERR] redis: pong not received")
	} else {
		log.Println("redis: pong received")
	}

	go func() {
		for {
			t1 := time.Now()
			err := stressTestRedis(client, 1000)
			dur := time.Since(t1)

			if err != nil {
				log.Println("[ERR]", err)
			} else {
				log.Printf("redis stress test: %s", dur.Round(time.Millisecond))
			}
			time.Sleep(10 * time.Second)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	log.Println("Exiting...")
}
