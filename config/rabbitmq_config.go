package config

import (
	"os"

	"github.com/joho/godotenv"
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
