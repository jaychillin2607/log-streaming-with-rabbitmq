package main

import (
	"encoding/json"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func receive(topicKey *string) {
	logger := newMultiLogger()

	// connection
	conn, err := amqp.Dial(RABBITMQ_URI)
	handleFailure("Error while connecting with RabbitMQ Broker", &err)
	defer conn.Close()

	// channel
	channel, err := conn.Channel()
	handleFailure("Error while creating channel", &err)
	defer channel.Close()

	// exchange
	err = channel.ExchangeDeclare(EXCHANGE_TOPIC, "topic", false, true, false, false, nil)
	handleFailure("Error while creating Exchange 'logs_topic'", &err)

	// queue
	queue, err := channel.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	handleFailure("Error while creating queue", &err)

	// bind
	for level := range strings.SplitSeq(strings.TrimSpace(strings.Trim(*topicKey, ".")), ".") {
		err = channel.QueueBind(
			queue.Name,
			createExchangeKey(&level),
			EXCHANGE_TOPIC,
			false,
			nil,
		)
		handleFailure("Error while binding Queue with the exchange", &err)

	}

	// consume
	channel.Qos(1, 0, false)
	msgChan, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleFailure("Error while consuming mesage from queue", &err)

	go func() {
		var jsonMsg logMessage
		for msg := range msgChan {

			json.Unmarshal(msg.Body, &jsonMsg)
			switch jsonMsg.Level {
			case INFO:
				logger.Info(jsonMsg.Msg)
			case DEBUG:
				logger.Debug(jsonMsg.Msg)
			case WARNING:
				logger.Warning(jsonMsg.Msg)
			case ERROR:
				logger.Error(jsonMsg.Msg)
			default:
				msg.Reject(false)
				handleFailure("Invalid log level", new(error))
			}
			msg.Ack(false)
		}

	}()

	var forever chan struct{}
	logger.Info(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
