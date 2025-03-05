package main

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

type Stresser struct {
	client *redis.Client

	count    int64 // Кол-во выполненных запросов
	channels []chan struct{}
}

func NewStresser(client *redis.Client) *Stresser {
	return &Stresser{
		client: client,
	}
}

func (s *Stresser) Inc() {
	done := make(chan struct{})
	s.channels = append(s.channels, done)
	go s.process(done)
}

func (s *Stresser) process(done chan struct{}) {
	for {
		err := s.oneGet()
		if err != nil {
			log.Printf("process died: %s", err)
			return
		}

		select {
		case <-done:
			return
		default:
		}
	}
}

func (s *Stresser) oneGet() error {
	err := s.client.Get(context.Background(), "key").Err()
	if err != nil {
		return err
	}
	atomic.AddInt64(&s.count, 1)
	return nil
}

func (s *Stresser) Dec() {
	if len(s.channels) <= 1 {
		return
	}
	done := s.channels[len(s.channels)-1]
	close(done)
	s.channels = s.channels[0 : len(s.channels)-1]
}

func (s *Stresser) Count() int {
	c := atomic.LoadInt64(&s.count)
	return int(c)
}

func (s *Stresser) Processes() int {
	return len(s.channels)
}
