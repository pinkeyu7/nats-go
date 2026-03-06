package service

import (
	"context"
	"nats-go/agent/config"
	"nats-go/pkg/jetstream"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	js "github.com/nats-io/nats.go/jetstream"
)

type KeepAliveServiceInterface interface {
	Start(ctx context.Context)
}

type KeepAliveService struct {
	js js.JetStream
}

func NewKeepAliveService(jsContext js.JetStream) KeepAliveServiceInterface {
	return &KeepAliveService{
		js: jsContext,
	}
}

func (s *KeepAliveService) Start(ctx context.Context) {
	// Send a keep-alive message every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Publish a keep-alive message to JetStream
			logger.Info("Sending keep-alive message...")
			agentID := config.GetAgentID()
			subject := jetstream.SubjectHeartbeatPrefix + agentID

			publishCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, err := s.js.Publish(publishCtx, subject, []byte(agentID))
			cancel()

			if err != nil {
				logger.Errorf("Error publishing keep-alive: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
