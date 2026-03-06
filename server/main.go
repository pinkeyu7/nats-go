package main

import (
	"flag"
	"log"
	"nats-go/server/api"
	"nats-go/server/config"
	"nats-go/server/route"

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
		log.Printf("Error initializing config: %v", err)
		return
	}

	// init api env
	err = api.InitEnv()
	if err != nil {
		log.Printf("Error initializing env: %v", err)
		return
	}

	// init gin router
	r := route.Init()

	// start server
	err = r.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
