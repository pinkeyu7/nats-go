package service

import (
	"context"
	"log"
	"nats-go/agent/config"
	"nats-go/pkg/topic"
	"time"

	"github.com/nats-io/nats.go"
)

type KeepAliveServiceInterface interface {
	Start(ctx context.Context)
}

type KeepAliveService struct {
	nc *nats.Conn
}

func NewKeepAliveService(nc *nats.Conn) KeepAliveServiceInterface {
	return &KeepAliveService{
		nc: nc,
	}
}

func (s *KeepAliveService) Start(ctx context.Context) {
	// Send a keep-alive message every 30 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Publish a keep-alive message to the "keep-alive" subject
			if err := s.nc.Publish(topic.TopicAgentKeepAlive, []byte(config.GetAgentID())); err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
