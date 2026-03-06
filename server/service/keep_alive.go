package service

import (
	"nats-go/pkg/topic"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/nats-io/nats.go"
)

type KeepAliveServiceInterface interface {
	Close() error
}

type KeepAliveService struct {
	nc  *nats.Conn
	sub *nats.Subscription
}

func NewKeepAliveService(nc *nats.Conn) KeepAliveServiceInterface {
	ks := &KeepAliveService{
		nc: nc,
	}

	var err error
	// Subscribe to tasks subject
	ks.sub, err = nc.Subscribe(topic.TopicAgentKeepAlive, func(msg *nats.Msg) {
		logger.Infof("Received keep-alive message: %s", string(msg.Data))
	})
	if err != nil {
		logger.Fatalf("Error subscribing to tasks: %v", err)
	}

	logger.Info("Agent is listening for tasks on 'tasks' subject...")

	return ks
}

func (s *KeepAliveService) Close() error {
	if err := s.sub.Unsubscribe(); err != nil {
		logger.Infof("Error unsubscribing: %v", err)
		return err
	}
	return nil
}
