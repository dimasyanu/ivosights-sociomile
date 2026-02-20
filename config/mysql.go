package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type MysqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func NewMysqlConfig(path ...string) *MysqlConfig {
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

	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		panic("Invalid MYSQL_PORT: " + err.Error())
	}
	return &MysqlConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     port,
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	}
}
