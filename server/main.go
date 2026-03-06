package main

import (
	"flag"
	"nats-go/server/api"
	"nats-go/server/config"
	"nats-go/server/route"
	"nats-go/server/service"

	"github.com/bytedance/gopkg/util/logger"
	_ "github.com/joho/godotenv/autoload"
)

var port string

func main() {
	// init http port
	flag.StringVar(&port, "port", "8080", "Initial port number")
	flag.Parse()

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

	apiEnv := api.GetEnv()
	// start keep-alive service
	ks := service.NewKeepAliveService(apiEnv.GetNC())
	defer func() {
		_ = ks.Close()
	}()

	// init gin router
	r := route.Init()

	// start server
	err = r.Run(":" + port)
	if err != nil {
		logger.Error(err)
	}
}
