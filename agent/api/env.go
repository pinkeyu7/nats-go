package api

import (
	"nats-go/agent/config"

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

	return nil
}

func (e *Env) Close() {
	if e.nc != nil {
		e.nc.Close()
	}
}
