package config

import (
	"os"
)

type JwtConfig struct {
	SecretKey string
}

func NewJwtConfig(path ...string) *JwtConfig {
	LoadEnvFile(path)
	return LoadJwtConfig()
}

func LoadJwtConfig() *JwtConfig {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		panic("JWT_SECRET_KEY environment variable is not set")
	}

	return &JwtConfig{
		SecretKey: secretKey,
	}
}
