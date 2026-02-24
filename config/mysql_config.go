package config

import (
	"fmt"
	"os"
	"strconv"
)

type MysqlConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Dsn      string
}

func NewMysqlConfig(path ...string) *MysqlConfig {
	LoadEnvFile(path)
	return LoadMysqlConfig()
}

func LoadMysqlConfig() *MysqlConfig {
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		panic("Invalid MYSQL_PORT: " + err.Error())
	}
	c := &MysqlConfig{
		Driver:   "mysql",
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     port,
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	}
	c.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
	return c
}
