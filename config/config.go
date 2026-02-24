package config

import (
	"github.com/joho/godotenv"
)

const EnvPath = ".env"

type Config struct {
	Http     *RestConfig
	MySQL    *MysqlConfig
	Jwt      *JwtConfig
	RabbitMQ *RabbitMQConfig
}

func NewConfig(path ...string) *Config {
	LoadEnvFile(path)
	return &Config{
		Http:     LoadRestConfig(),
		MySQL:    LoadMysqlConfig(),
		Jwt:      LoadJwtConfig(),
		RabbitMQ: LoadRabbitMQConfig(),
	}
}

func LoadEnvFile(path []string) string {
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
	return envPath
}
