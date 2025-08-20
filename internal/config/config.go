package config

import (
	"errors"
	"os"
)

const RABBITMQ_URI string = "amqp://guest:guest@localhost:5672/"
const EXCHANGE_TOPIC string = "logs_topic"

var DEFAULT_LOG_LEVEL string = getEnvOrDefault("DEFAULT_LOG_LEVEL", "INFO")
var DEFAULT_TOPIC string = getEnvOrDefault("DEFAULT_TOPIC", "*.DEBUG.*")
var EMPTY_ERROR error = errors.New("")

func getEnvOrDefault(envVar, defaultVar string) string {
	if val := os.Getenv(envVar); val != "" {
		return val
	}
	return defaultVar
}
