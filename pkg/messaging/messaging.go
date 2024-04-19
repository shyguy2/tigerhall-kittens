package messaging

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// MessageBroker represents the messaging service using RabbitMQ.
type MessageBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewMessageBroker creates a new MessageBroker instance.
func NewMessageBroker(amqpURL, queueName string) (*MessageBroker, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %v", err)
	}

	queue, err := channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare RabbitMQ queue: %v", err)
	}

	return &MessageBroker{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

// PublishMessage publishes a message to the RabbitMQ queue.
func (mb *MessageBroker) PublishMessage(message []byte) error {
	err := mb.channel.Publish(
		"",            // exchange
		mb.queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to RabbitMQ: %v", err)
	}

	return nil
}

// ConsumeMessages starts consuming messages from the RabbitMQ queue.
// It takes a message processing function as an argument to handle each message received.
func (mb *MessageBroker) ConsumeMessages(processMessage func([]byte) error) {
	msgs, err := mb.channel.Consume(
		mb.queue.Name, // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		err := processMessage(msg.Body)
		if err != nil {
			log.Printf("failed to process message: %v", err)
			// Requeue the message to be processed later
			msg.Nack(false, true)
		} else {
			// Acknowledge the successful processing of the message
			msg.Ack(false)
		}
	}
}

// Close closes the connection and channel to the RabbitMQ broker.
func (mb *MessageBroker) Close() {
	mb.channel.Close()
	mb.conn.Close()
}

func ProcessMessage(message []byte) error {
	log.Printf("Following list of email are sent: %s\n", message)
	return nil
}
