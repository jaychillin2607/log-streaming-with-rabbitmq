package internal

import (
	"log"

	"github.com/jaychillin2607/log-streaming-with-rabbitmq/internal/config"
	"github.com/jaychillin2607/log-streaming-with-rabbitmq/pkg/types"
)

func HandleFailure(errMsg string, err *error) {
	if *err != nil {
		log.Fatalf("Message: %s || Error: %v", errMsg, *err)
	}
}

func CastStrToLogLevel(level *string) types.LogLevel {
	switch *level {
	case "INFO":
		return types.INFO
	case "DEBUG":
		return types.DEBUG
	case "WARNING":
		return types.WARNING
	case "ERROR":
		return types.ERROR
	default:
		HandleFailure("Invalid log level provided. Valid levels [INFO, DEBUG, WARNING, ERROR]", &config.EMPTY_ERROR)
		return -1
	}
}

func CreateExchangeKey(level *string) string {
	switch *level {
	case "INFO":
		return "INFO.*.*"
	case "DEBUG":
		return "*.DEBUG.*"
	case "WARNING":
		return "*.*.WARNING"
	case "ERROR":
		return "ERROR.*"
	default:
		HandleFailure("Invalid log level provided. Valid levels [INFO, DEBUG, WARNING, ERROR]", &config.EMPTY_ERROR)
		return config.DEFAULT_TOPIC
	}
}
