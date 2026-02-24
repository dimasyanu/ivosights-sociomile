package config

import (
	"os"
	"strconv"
)

type RestConfig struct {
	Port uint16
}

func NewRestConfig(path ...string) *RestConfig {
	LoadEnvFile(path)
	return LoadRestConfig()
}

func LoadRestConfig() *RestConfig {
	port, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		panic("Invalid HTTP_PORT: " + err.Error())
	}

	return &RestConfig{
		Port: uint16(port),
	}
}
