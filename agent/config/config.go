package config

import "os"

func Init() error {
	return nil
}

func GetAgentID() string {
	id := os.Getenv("AGENT_ID")
	if id == "" {
		return "agent-001"
	}
	return id
}

func GetNatsURL() string {
	url := os.Getenv("NATS_URL")
	if url == "" {
		return "nats://localhost:4222"
	}
	return url
}
