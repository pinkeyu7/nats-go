package api

import (
	"nats-go/agent/config"

	"github.com/nats-io/nats.go"
)

type Env struct {
	nc *nats.Conn
}

var env = &Env{}

func GetEnv() *Env {
	return env
}

func (e *Env) GetNC() *nats.Conn {
	return e.nc
}

func InitEnv() error {
	var err error

	// Connect to NATS server
	env.nc, err = nats.Connect(config.GetNatsURL())
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
