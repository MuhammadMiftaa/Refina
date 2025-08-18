package queue

import (
	"fmt"

	"server/config/env"
	"server/config/log"

	"github.com/rabbitmq/amqp091-go"
)

var connection *amqp091.Connection

func SetupRabbitMQ() {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", env.Cfg.RabbitMQ.RMQUser, env.Cfg.RabbitMQ.RMQPassword, env.Cfg.RabbitMQ.RMQHost, env.Cfg.RabbitMQ.RMQPort, env.Cfg.RabbitMQ.RMQVirtualHost)

	conn, err := amqp091.Dial(connectionString)
	if err != nil {
		log.Log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	connection = conn
}

func Close() {
	if connection != nil {
		if err := connection.Close(); err != nil {
			log.Error("Failed to close RabbitMQ connection: " + err.Error())
		}
	}
}

func GetChannel() (*amqp091.Channel, error) {
	if connection == nil {
		return nil, fmt.Errorf("RabbitMQ connection is not initialized")
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return channel, nil
}
