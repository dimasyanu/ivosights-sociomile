package infra

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ConversationQueue = "conversation_queue"

type RabbitMQConfig struct {
	URL         string
	Exchange    string
	Queue       string
	RoutingKey  string
	ContentType string
	Username    string
	Password    string
}

func NewRabbitMQConfig(path ...string) *RabbitMQConfig {
	var envPath string
	if len(path) > 0 {
		envPath = path[0]
	} else {
		envPath = ".env"
	}

	if len(envPath) > 0 {
		if err := godotenv.Load(envPath); err != nil {
			panic("Error loading configuration: " + err.Error())
		}
	}
	c := &RabbitMQConfig{
		URL:         os.Getenv("RABBITMQ_URL"),
		Exchange:    os.Getenv("RABBITMQ_EXCHANGE"),
		Queue:       os.Getenv("RABBITMQ_QUEUE"),
		RoutingKey:  os.Getenv("RABBITMQ_ROUTING_KEY"),
		Username:    os.Getenv("RABBITMQ_DEFAULT_USER"),
		Password:    os.Getenv("RABBITMQ_DEFAULT_PASS"),
		ContentType: "application/json",
	}
	return c
}

type RabbitMQClient struct {
	config *RabbitMQConfig
	conn   *amqp.Connection
}

func NewRabbitMQClient(config *RabbitMQConfig) (QueueClient, error) {
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

func (c *RabbitMQClient) PublishMessage(queue string, message []byte) error {
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

	defer func() {
		err = msg.Reject(true)
		if err != nil {
			fmt.Printf("Error rejecting message: %v\n", err)
		}
	}()
	return [][]byte{msg.Body}
}

func (c *RabbitMQClient) Close() error {
	c.conn.Close()
	return nil
}
