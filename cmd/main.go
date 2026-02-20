package main

import (
	"database/sql"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest"

	_ "github.com/dimasyanu/ivosights-sociomile/docs"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := ConnectToDatabase()
	defer db.Close()

	restConfig := config.NewRestConfig(config.EnvPath)
	api := rest.NewRestApi(restConfig, db)
	api.Start()
}

func ConnectToDatabase() *sql.DB {
	mysqlCfg := config.NewMysqlConfig(config.EnvPath)
	db, err := sql.Open(mysqlCfg.Driver, mysqlCfg.Dsn)
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		panic("Error pinging database: " + err.Error())
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}
