package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

// fillRedis наполнит Редис данными (n ключей).
func fillRedis(client *redis.Client, n int) error {
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%d", i)
		val := fmt.Sprintf("%d", i)

		err := client.Set(context.Background(), key, val, 0).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

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

type Stresser struct {
	client *redis.Client

	counter  int64 // Кол-во выполненных запросов
	channels []chan struct{}
}

func NewStresser(client *redis.Client) *Stresser {
	return &Stresser{
		client: client,
	}
}

func (s *Stresser) AddProccess() {
	done := make(chan struct{})
	s.channels = append(s.channels, done)
	go func() {
		for {
			err := s.client.Get(context.Background(), "key").Err()
			if err != nil {
				log.Printf("process died: %s", err)
				return
			}

			atomic.AddInt64(&s.counter, 1)

			select {
			case <-done:
				return
			default:
			}
		}
	}()
}

func (s *Stresser) SubProccess() {
	done := s.channels[len(s.channels)-1]
	close(done)
	s.channels = s.channels[0 : len(s.channels)-1]
}

func (s *Stresser) GetCount() int {
	c := atomic.LoadInt64(&s.counter)
	return int(c)
}

func main() {
	// Создадим подключение к Редис
	//
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Printf("[ERR] redis: ping: %s", err)
	} else {
		log.Println("redis: pong received")
	}

	// Наполним Редис первоначальными данными
	err = fillRedis(client, 1000)
	if err != nil {
		log.Printf("[ERR] redis: наполнение данными: %s", err)
	} else {
		log.Println("redis: pong received")
	}

	s := NewStresser(client)
	s.AddProccess()

	go func() {
		for {
			s.GetCount()

			time.Sleep(5 * time.Second)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	log.Println("Exiting...")
}
