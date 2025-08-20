package main

import (
	"flag"
	"log"
	"os"

	"github.com/jaychillin2607/log-streaming-with-rabbitmq/internal/config"
	"github.com/jaychillin2607/log-streaming-with-rabbitmq/internal/rabbitmq"
)

func main() {

	sendCommand := flag.NewFlagSet("send", flag.ExitOnError)
	logLevel := sendCommand.String("level", config.DEFAULT_LOG_LEVEL, "sets the level of the log message")
	logMsg := sendCommand.String("msg", "", "sets the log message")
	receiveCommand := flag.NewFlagSet("receive", flag.ExitOnError)
	topicKey := receiveCommand.String("topicKey", config.DEFAULT_TOPIC, "sets the topic binding key")

	switch os.Args[1] {
	case "send":
		sendCommand.Parse(os.Args[2:])
		rabbitmq.Send(logLevel, logMsg)

	case "receive":
		receiveCommand.Parse(os.Args[2:])
		rabbitmq.Receive(topicKey)
	default:
		log.Fatalf("invalid subcommand: %s", os.Args[1])
	}
}
