package config

import (
	"os"
)

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
	LoadEnvFile(path)
	return LoadRabbitMQConfig()
}

func LoadRabbitMQConfig() *RabbitMQConfig {
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
