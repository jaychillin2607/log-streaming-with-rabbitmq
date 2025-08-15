package main

import (
	"errors"
	"flag"
	"log"
	"os"
)

const RABBITMQ_URI string = "amqp://guest:guest@localhost:5672/"
const EXCHANGE_TOPIC string = "logs_topic"

var DEFAULT_LOG_LEVEL string = getEnvOrDefault("DEFAULT_LOG_LEVEL", "INFO")
var DEFAULT_TOPIC string = getEnvOrDefault("DEFAULT_TOPIC", "*.DEBUG.*")
var EMPTY_ERROR error = errors.New("")

func main() {

	sendCommand := flag.NewFlagSet("send", flag.ExitOnError)
	logLevel := sendCommand.String("level", DEFAULT_LOG_LEVEL, "sets the level of the log message")
	logMsg := sendCommand.String("msg", "", "sets the log message")
	receiveCommand := flag.NewFlagSet("receive", flag.ExitOnError)
	topicKey := receiveCommand.String("topicKey", DEFAULT_TOPIC, "sets the topic binding key")

	switch os.Args[1] {
	case "send":
		sendCommand.Parse(os.Args[2:])
		send(logLevel, logMsg)

	case "receive":
		receiveCommand.Parse(os.Args[2:])
		receive(topicKey)
	default:
		log.Fatalf("invalid subcommand: %s", os.Args[1])
	}
}
