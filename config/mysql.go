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

func NewMysqlConfig() (*MysqlConfig, error) {
	if err := godotenv.Load(); err != nil {
		panic("Error loading configuration: " + err.Error())
	}

	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		return nil, err
	}
	return &MysqlConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     port,
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	}, nil
}
