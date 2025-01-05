package messaging

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
	QueueName  string
}

func InitRabbitMQ(queueName string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial("amqp://user:password@rabbitmq:5672/") // "amqp://user:password@rabbitmq:5672/"
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close() // Ensure connection is closed if channel creation fails
		return nil, fmt.Errorf("failed to create RabbitMQ channel: %w", err)
	}

	// Declare a queue
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	log.Printf("RabbitMQ initialized with queue: %s", queueName)
	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		QueueName:  queueName,
	}, nil
}

func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
	log.Println("RabbitMQ connection and channel closed.")
}
