package main

import (
	"context"
	"nats-go/agent/api"
	"nats-go/agent/config"
	"nats-go/agent/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/bytedance/gopkg/util/logger"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger.Info("Agent is starting...")

	// init config
	err := config.Init()
	if err != nil {
		logger.Infof("Error initializing config: %v", err)
		return
	}

	// init api env
	err = api.InitEnv()
	if err != nil {
		logger.Infof("Error initializing env: %v", err)
		return
	}
	defer api.GetEnv().Close()

	// setup global context for keep-alive service
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	apiEnv := api.GetEnv()

	// start keep-alive service
	ks := service.NewKeepAliveService(apiEnv.GetJS())
	go ks.Start(ctx)

	// start task service
	ts := service.NewTaskService(apiEnv.GetNC(), apiEnv.GetJS())
	defer func() {
		ts.Close()
	}()

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Agent is shutting down...")
}
