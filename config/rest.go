package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type RestConfig struct {
	Port uint16
}

func NewRestConfig() *RestConfig {
	if err := godotenv.Load(); err != nil {
		panic("Error loading configuration: " + err.Error())
	}

	port, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		panic("Invalid HTTP_PORT: " + err.Error())
	}

	return &RestConfig{
		Port: uint16(port),
	}
}
