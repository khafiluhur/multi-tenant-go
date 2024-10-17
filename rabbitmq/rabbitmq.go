package rabbitmq

import (
	"jatis_mobile_api/logs"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var (
	conn    *amqp091.Connection
	channel *amqp091.Channel
	logger  = logs.SetupLogger()
)

func ConnectRabbitMQ(url string) error {
	var err error
	conn, err = amqp091.Dial(url)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to connect to RabbitMQ", struct{ RabbitMQURL string }{RabbitMQURL: url})
		return err
	}
	logs.LogWithFields(logger, logrus.InfoLevel, "Connected to RabbitMQ", struct{ RabbitMQURL string }{RabbitMQURL: url})

	channel, err = conn.Channel()
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to open a channel", struct{ Error error }{Error: err})
		return err
	}
	logs.LogWithFields(logger, logrus.InfoLevel, "Channel opened successfully", struct{}{})
	return nil
}

// type exchange direct
func PublishMessage(exchangeName, routingKey string, body []byte) error {
	if channel == nil {
		err := amqp091.ErrClosed
		logs.LogWithFields(logger, logrus.ErrorLevel, "Channel is not available", struct{ Error error }{Error: err})
		return amqp091.ErrClosed
	}

	err := channel.Publish(
		exchangeName, // The exchange to publish to
		routingKey,   // The routing key used for routing the message
		false,        // Mandatory
		false,        // Immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to publish message", struct{ RoutingKey string }{RoutingKey: routingKey})
		return err
	}

	logs.LogWithFields(logger, logrus.InfoLevel, "Message published successfully", struct {
		RoutingKey string
		Body       string
	}{RoutingKey: routingKey, Body: string(body)})
	return nil
}

func DeclareQueue(queueName string) error {
	_, err := channel.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		amqp091.Table{
			"x-queue-type": "quorum", // Use quorum queue
		},
	)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to declare queue", struct{ QueueName string }{QueueName: queueName})
		return err
	}

	logs.LogWithFields(logger, logrus.InfoLevel, "Queue declared successfully", struct{ QueueName string }{QueueName: queueName})
	return nil
}

func BindQueue(queueName, exchangeName, routingKey string) error {
	err := channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to bind queue to exchange", struct {
			QueueName    string
			ExchangeName string
			RoutingKey   string
		}{
			QueueName:    queueName,
			ExchangeName: exchangeName,
			RoutingKey:   routingKey,
		})
		return err
	}

	logs.LogWithFields(logger, logrus.InfoLevel, "Queue bound to exchange successfully", struct {
		QueueName    string
		ExchangeName string
		RoutingKey   string
	}{
		QueueName:    queueName,
		ExchangeName: exchangeName,
		RoutingKey:   routingKey,
	})

	return nil
}

func ConsumeMessages(queueName string) error {
	msgs, err := channel.Consume(
		queueName,
		"",
		false, // Auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to register consumer", struct{ QueueName string }{QueueName: queueName})
		return err
	}

	logs.LogWithFields(logger, logrus.InfoLevel, "Consumer started successfully", struct{ QueueName string }{QueueName: queueName})

	go func() {
		for msg := range msgs {
			processMessage(msg)
		}
	}()

	return nil
}

func processMessage(msg amqp091.Delivery) {
	logs.LogWithFields(logger, logrus.InfoLevel, "Processing message", struct{ Body string }{Body: string(msg.Body)})

	if err := msg.Ack(false); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to acknowledge message", struct{ Error error }{Error: err})
	} else {
		logs.LogWithFields(logger, logrus.InfoLevel, "Message acknowledged", struct{ DeliveryTag uint64 }{DeliveryTag: msg.DeliveryTag})
	}
}

func GetChannel() (*amqp091.Channel, error) {
	if channel == nil {
		err := amqp091.ErrClosed
		logs.LogWithFields(logger, logrus.ErrorLevel, "Channel is not available", struct{ Error error }{Error: err})
		return nil, amqp091.ErrClosed
	}
	logs.LogWithFields(logger, logrus.InfoLevel, "Channel retrieved successfully", struct{}{})
	return channel, nil
}

func Close() {
	if channel != nil {
		if err := channel.Close(); err != nil {
			logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to close channel", struct{ Error error }{Error: err})
		} else {
			logs.LogWithFields(logger, logrus.InfoLevel, "Channel closed successfully", struct{ Error error }{Error: err})
		}
	}
	if conn != nil {
		if err := conn.Close(); err != nil {
			logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to close connection", struct{ Error error }{Error: err})
		} else {
			logs.LogWithFields(logger, logrus.InfoLevel, "Connection closed successfully", struct{ Error error }{Error: err})
		}
	}
}

func IsClosed() bool {
	return conn == nil || conn.IsClosed()
}
