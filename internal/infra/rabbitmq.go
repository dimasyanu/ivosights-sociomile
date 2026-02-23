package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ConversationQueue = "conversation_queue"

type RabbitMQClient struct {
	config *config.RabbitMQConfig
	conn   *amqp.Connection
}

func NewRabbitMQClient(config *config.RabbitMQConfig) (QueueClient, error) {
	// Establish connection to RabbitMQ server
	url := fmt.Sprintf(config.URL, config.Username, config.Password)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	// Set up channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare(
		config.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Set up exchange and routing key
	err = ch.ExchangeDeclare(
		config.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Bind the queue to the exchange with the routing key
	err = ch.QueueBind(
		config.Queue,
		config.RoutingKey,
		config.Exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{config: config, conn: conn}, nil
}

func (c *RabbitMQClient) Publish(queue string, message []byte) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx,
		c.config.Exchange,
		c.config.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: c.config.ContentType,
			Body:        message,
		},
	)
	return err
}

func (c *RabbitMQClient) Consume(queue string) (<-chan []byte, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	messageChan := make(chan []byte)
	go func() {
		defer ch.Close()
		for msg := range msgs {
			messageChan <- msg.Body
		}
		close(messageChan)
	}()

	return messageChan, nil
}

func (c *RabbitMQClient) GetPublishedMessages() [][]byte {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil
	}
	defer ch.Close()

	msg, ok, err := ch.Get(c.config.Queue, true)
	if err != nil {
		fmt.Printf("Error getting message: %v\n", err)
		return nil
	}
	if !ok {
		fmt.Printf("No message available\n")
		return nil
	}

	fmt.Printf("Message: %s\n", string(msg.Body))

	return [][]byte{msg.Body}
}

func (c *RabbitMQClient) PeekPublishedMessages() [][]byte {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil
	}
	defer ch.Close()

	msg, ok, err := ch.Get(c.config.Queue, true)
	if err != nil {
		fmt.Printf("Error getting message: %v\n", err)
		return nil
	}
	if !ok {
		fmt.Printf("No message available\n")
		return nil
	}

	fmt.Printf("Message: %s\n", string(msg.Body))

	defer func() {
		err = msg.Reject(true)
		if err != nil {
			fmt.Printf("Error rejecting message: %v\n", err)
		}
	}()
	return [][]byte{msg.Body}
}

func (c *RabbitMQClient) Clear() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueuePurge(c.config.Queue, false)
	if err != nil {
		return err
	}
	return nil
}

func (c *RabbitMQClient) Close() error {
	c.conn.Close()
	return nil
}
