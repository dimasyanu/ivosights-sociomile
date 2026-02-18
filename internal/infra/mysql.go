package infra

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlDatabase struct {
	Db *sql.DB
}

var once sync.Once

func NewMySQLDatabase(conf *config.Config) Database {
	var database *mysqlDatabase
	once.Do(func() {
		dsn := ""
		gormDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to database: %v", err.Error()))
		}
		log.Printf("Successfully connected to MySQL database.")
		db, err := gormDb.DB()
		database = &mysqlDatabase{Db: db}
	})
	return database
}

func (d *mysqlDatabase) GetDb() (*sql.DB, error) {
	return d.Db, nil
}

func (d *mysqlDatabase) Close() error {
	return d.Db.Close()
}
