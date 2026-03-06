package service

import (
	"context"
	"nats-go/pkg/jetstream"

	"github.com/bytedance/gopkg/util/logger"
	js "github.com/nats-io/nats.go/jetstream"
)

type KeepAliveServiceInterface interface {
	Close() error
}

type KeepAliveService struct {
	js       js.JetStream
	consumer js.Consumer
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewKeepAliveService(jsContext js.JetStream) KeepAliveServiceInterface {
	ks := &KeepAliveService{
		js: jsContext,
	}

	ks.ctx, ks.cancel = context.WithCancel(context.Background())

	var err error
	// Get the consumer
	ks.consumer, err = jsContext.Consumer(ks.ctx, jetstream.StreamHeartbeat, jetstream.ConsumerHeartbeatMonitor)
	if err != nil {
		logger.Fatalf("Error getting heartbeat consumer: %v", err)
	}

	// Start consuming messages
	go ks.consumeHeartbeats()

	logger.Info("Server is monitoring agent heartbeats via JetStream...")

	return ks
}

func (s *KeepAliveService) consumeHeartbeats() {
	// Consume heartbeat messages from the stream
	cons, err := s.consumer.Consume(func(msg js.Msg) {
		agentID := string(msg.Data())
		logger.Infof("Received keep-alive from agent: %s", agentID)

		// Acknowledge successful processing
		if err := msg.Ack(); err != nil {
			logger.Errorf("Error acknowledging message: %v", err)
			return
		}
	})

	if err != nil {
		logger.Fatalf("Error consuming heartbeat messages: %v", err)
	}

	// Wait for context cancellation
	<-s.ctx.Done()
	cons.Stop()
}

func (s *KeepAliveService) Close() error {
	logger.Info("Closing KeepAliveService...")
	s.cancel()
	return nil
}
