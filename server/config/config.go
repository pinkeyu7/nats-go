package config

import (
	"os"
	"strings"
)

const (
	EnvProduction  = "production"
	EnvStaging     = "staging"
	EnvDevelopment = "development"
	EnvLocalhost   = "localhost"
)

var EnvShortName = map[string]string{
	EnvProduction:  "prod",
	EnvStaging:     "stag",
	EnvDevelopment: "dev",
	EnvLocalhost:   "local",
}

func Init() error {
	return nil
}

func GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}

func GetCorsRule(origin string) bool {
	switch GetEnvironment() {
	case EnvLocalhost:
		return true
	case EnvDevelopment:
		return origin == "https://sample-development.website.com" || strings.Contains(origin, "http://localhost")
	case EnvStaging:
		return origin == "https://sample-staging.website.com"
	case EnvProduction:
		return origin == "https://sample.website.com"
	default:
		return true
	}
}

func GetDBPath() string {
	return os.Getenv("DB_PATH")
}

func GetNatsURL() string {
	url := os.Getenv("NATS_URL")
	if url == "" {
		return "nats://localhost:4222"
	}
	return url
}
