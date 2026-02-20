package config

import (
	"os"

	"github.com/joho/godotenv"
)

type JwtConfig struct {
	SecretKey string
}

func NewJwtConfig(path ...string) *JwtConfig {
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

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		panic("JWT_SECRET_KEY environment variable is not set")
	}

	return &JwtConfig{
		SecretKey: secretKey,
	}
}
