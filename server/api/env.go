package api

type Env struct {
}

var env = &Env{}

func GetEnv() *Env {
	return env
}

func InitEnv() error {
	return nil
}
