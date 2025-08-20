package rabbitmq

import (
	"encoding/json"
	"strings"

	"github.com/jaychillin2607/log-streaming-with-rabbitmq/internal"
	"github.com/jaychillin2607/log-streaming-with-rabbitmq/internal/config"
	"github.com/jaychillin2607/log-streaming-with-rabbitmq/internal/logger"
	"github.com/jaychillin2607/log-streaming-with-rabbitmq/pkg/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Receive(topicKey *string) {
	logger := logger.NewMultiLogger()

	// connection
	conn, err := amqp.Dial(config.RABBITMQ_URI)
	internal.HandleFailure("Error while connecting with RabbitMQ Broker", &err)
	defer conn.Close()

	// channel
	channel, err := conn.Channel()
	internal.HandleFailure("Error while creating channel", &err)
	defer channel.Close()

	// exchange
	err = channel.ExchangeDeclare(config.EXCHANGE_TOPIC, "topic", false, true, false, false, nil)
	internal.HandleFailure("Error while creating Exchange 'logs_topic'", &err)

	// queue
	queue, err := channel.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	internal.HandleFailure("Error while creating queue", &err)

	// bind
	for level := range strings.SplitSeq(strings.TrimSpace(strings.Trim(*topicKey, ".")), ".") {
		err = channel.QueueBind(
			queue.Name,
			internal.CreateExchangeKey(&level),
			config.EXCHANGE_TOPIC,
			false,
			nil,
		)
		internal.HandleFailure("Error while binding Queue with the exchange", &err)

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
	internal.HandleFailure("Error while consuming mesage from queue", &err)

	go func() {
		var jsonMsg types.LogMessage
		for msg := range msgChan {

			json.Unmarshal(msg.Body, &jsonMsg)
			switch jsonMsg.Level {
			case types.INFO:
				logger.Info(jsonMsg.Msg)
			case types.DEBUG:
				logger.Debug(jsonMsg.Msg)
			case types.WARNING:
				logger.Warning(jsonMsg.Msg)
			case types.ERROR:
				logger.Error(jsonMsg.Msg)
			default:
				msg.Reject(false)
				internal.HandleFailure("Invalid log level", new(error))
			}
			msg.Ack(false)
		}

	}()

	var forever chan struct{}
	logger.Info(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func Send(level, msg *string) {
	// connection
	conn, err := amqp.Dial(config.RABBITMQ_URI)
	internal.HandleFailure("Error while connecting with RabbitMQ Broker", &err)
	defer conn.Close()

	// channel
	channel, err := conn.Channel()
	internal.HandleFailure("Error while creating channel", &err)
	defer channel.Close()

	// exchange
	err = channel.ExchangeDeclare(
		config.EXCHANGE_TOPIC,
		"topic",
		false,
		true,
		false,
		false,
		nil,
	)
	internal.HandleFailure("Error while creating Exchange 'logs_topic'", &err)

	// publish
	message, err := json.Marshal(*logger.NewLogMessage(level, msg))
	internal.HandleFailure("Error while marshalling message", &err)

	channel.Publish(
		config.EXCHANGE_TOPIC,
		internal.CreateExchangeKey(level),
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
