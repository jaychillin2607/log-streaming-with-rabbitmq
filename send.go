package main

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func send(level, msg *string) {
	// connection
	conn, err := amqp.Dial(RABBITMQ_URI)
	handleFailure("Error while connecting with RabbitMQ Broker", &err)
	defer conn.Close()

	// channel
	channel, err := conn.Channel()
	handleFailure("Error while creating channel", &err)
	defer channel.Close()

	// exchange
	err = channel.ExchangeDeclare(
		createExchangeKey(level),
		"topic",
		false,
		true,
		false,
		false,
		nil,
	)
	handleFailure("Error while creating Exchange 'logs_topic'", &err)

	// publish
	message, err := json.Marshal(*NewLogMessage(level, msg))
	handleFailure("Error while marshalling message", &err)

	channel.Publish(
		EXCHANGE_TOPIC,
		createExchangeKey(level),
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
