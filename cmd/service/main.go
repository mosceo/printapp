package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// stressTestRedis отправит n кол-во set запросов в Редис.
// func stressTestRedis(client *redis.Client, n int) error {
// 	for i := 0; i < n; i++ {
// 		err := client.Set(context.Background(), "key", "value", 0).Err()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func work() error {
	client, err := ConnectRedis()
	if err != nil {
		return fmt.Errorf("redis connect: %w", err)
	}

	err = FillRedis(client, 1000)
	if err != nil {
		return fmt.Errorf("redis fill: %w", err)
	}

	// Создадим пока один процесс и через интервалы будем
	// увеличивать или уменьшать кол-во процессов
	//
	stress := NewStresser(client)
	stress.Inc()
	calc := NewCalculator()

	go func() {
		for {
			time.Sleep(5 * time.Second)

			count := stress.Count()
			rate, grown := calc.Set(count)

			log.Printf("скорость %v | процессов %d", math.Round(rate), stress.Processes())

			if grown {
				stress.Inc()
			} else {
				stress.Dec()
			}
		}
	}()

	return nil
}

func main() {
	err := work()
	if err != nil {
		log.Printf("[ERR] %s", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	log.Println("Exiting...")
}
