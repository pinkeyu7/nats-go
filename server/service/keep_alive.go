package service

import (
	"log"

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
	ks.sub, err = nc.Subscribe("agent.Keep_alive", func(msg *nats.Msg) {
		log.Printf("Received keep-alive message: %s", string(msg.Data))
	})
	if err != nil {
		log.Fatalf("Error subscribing to tasks: %v", err)
	}

	log.Println("Agent is listening for tasks on 'tasks' subject...")

	return ks
}

func (s *KeepAliveService) Close() error {
	if err := s.sub.Unsubscribe(); err != nil {
		log.Printf("Error unsubscribing: %v", err)
		return err
	}
	return nil
}
