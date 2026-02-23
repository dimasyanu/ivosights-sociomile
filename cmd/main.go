package main

import (
	"database/sql"
	"sync"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"

	_ "github.com/dimasyanu/ivosights-sociomile/docs"
	_ "github.com/go-sql-driver/mysql"
)

var once sync.Once

func main() {
	db := ConnectToDatabase()
	defer db.Close()

	mq := ConnectToMq()
	defer mq.Close()

	restConfig := config.NewRestConfig(config.EnvPath)
	api := rest.NewRestApi(restConfig, db, mq)
	api.Start()
}

func ConnectToDatabase() *sql.DB {
	mysqlCfg := config.NewMysqlConfig(config.EnvPath)

	var db *sql.DB
	var err error
	once.Do(func() {
		db, err = infra.NewMySQLDatabase(mysqlCfg)
		if err != nil {
			panic("Error connecting to database: " + err.Error())
		}
	})

	if err := db.Ping(); err != nil {
		panic("Error pinging database: " + err.Error())
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}

func ConnectToMq() infra.QueueClient {
	rabbitMqCfg := config.NewRabbitMQConfig(config.EnvPath)

	mq, err := infra.NewRabbitMQClient(rabbitMqCfg)
	if err != nil {
		panic("Error connecting to RabbitMQ: " + err.Error())
	}

	return mq
}
