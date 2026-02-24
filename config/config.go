package config

const EnvPath = ".env"

type Config struct {
	Http     *RestConfig
	Mysql    *MysqlConfig
	Jwt      *JwtConfig
	RabbitMq *RabbitMQConfig
}

func NewConfig() *Config {
	return &Config{
		Http:     NewRestConfig(EnvPath),
		Mysql:    NewMysqlConfig(EnvPath),
		Jwt:      NewJwtConfig(EnvPath),
		RabbitMq: NewRabbitMQConfig(EnvPath),
	}
}
