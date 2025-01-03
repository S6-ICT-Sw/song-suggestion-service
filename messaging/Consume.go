package messaging

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) Consume() (<-chan amqp091.Delivery, error) {
	msgs, err := r.Channel.Consume(
		r.QueueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consuming messages from queue %s: %w", r.QueueName, err)
	}
	return msgs, nil
}
