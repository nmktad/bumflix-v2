package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQQueue is a RabbitMQ implementation of the Queue interface
type RabbitMQQueue struct {
	Conn *amqp.Connection
}

// NewRabbitMQQueue creates a new RabbitMQQueue

func NewRabbitMQQueue(conn *amqp.Connection) *RabbitMQQueue {
	return &RabbitMQQueue{Conn: conn}
}

// Publish publishes a message to the queue

func (q *RabbitMQQueue) Publish(queueName string, body []byte) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	return ch.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
}

// Consume consumes a message from the queue
func (q *RabbitMQQueue) Consume(queueName string) (<-chan amqp.Delivery, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return ch.Consume(queueName, "", false, false, false, false, nil)
}
