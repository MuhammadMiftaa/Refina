package queue

import (
	"fmt"

	"server/config/env"

	"github.com/rabbitmq/amqp091-go"
)

var connection *amqp091.Connection

func SetupRabbitMQ() {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", env.Cfg.RabbitMQ.RMQUser, env.Cfg.RabbitMQ.RMQPassword, env.Cfg.RabbitMQ.RMQHost, env.Cfg.RabbitMQ.RMQPort, env.Cfg.RabbitMQ.RMQVirtualHost)

	conn, err := amqp091.Dial(connectionString)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %s\n", err)
		panic(err)
	}

	connection = conn
	fmt.Println("Connected to RabbitMQ successfully")
}

func Close() error {
	if connection != nil {
		if err := connection.Close(); err != nil {
			return fmt.Errorf("failed to close RabbitMQ connection: %w", err)
		}
		fmt.Println("RabbitMQ connection closed successfully")
	}
	return nil
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
