package api

import (
	"context"
	"nats-go/pkg/jetstream"
	"nats-go/server/config"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/nats-io/nats.go"
	js "github.com/nats-io/nats.go/jetstream"
)

type Env struct {
	nc *nats.Conn
	js js.JetStream
}

var env = &Env{}

func GetEnv() *Env {
	return env
}

func (e *Env) GetNC() *nats.Conn {
	return e.nc
}

func (e *Env) GetJS() js.JetStream {
	return e.js
}

func InitEnv() error {
	var err error

	// Connect to NATS server
	env.nc, err = nats.Connect(config.GetNatsURL())
	if err != nil {
		return err
	}

	// Initialize JetStream
	env.js, err = js.New(env.nc)
	if err != nil {
		return err
	}

	// Setup JetStream streams
	logger.Info("Setting up JetStream streams...")
	ctx := context.Background()
	if err := jetstream.SetupStreams(ctx, env.js); err != nil {
		logger.Errorf("Failed to setup streams: %v", err)
		return err
	}

	// Setup JetStream consumers
	logger.Info("Setting up JetStream consumers...")
	if err := jetstream.SetupConsumers(ctx, env.js); err != nil {
		logger.Errorf("Failed to setup consumers: %v", err)
		return err
	}

	logger.Info("JetStream initialized successfully")
	return nil
}

func (e *Env) Close() {
	if e.nc != nil {
		e.nc.Close()
	}
}
