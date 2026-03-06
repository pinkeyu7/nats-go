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

const (
	PathImageUpload         = "./upload/image"
	PathImageCompressUpload = "./upload/image_compress"

	ExternalPathImageUpload = "/static/images"
)

func Init() error {
	err := os.MkdirAll(PathImageUpload, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(PathImageCompressUpload, os.ModePerm)
	if err != nil {
		return err
	}

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
