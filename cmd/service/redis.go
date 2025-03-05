package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// FillRedis наполнит Редис данными (n ключей).
func FillRedis(client *redis.Client, n int) error {
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
