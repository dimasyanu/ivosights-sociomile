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

func NewMySQLDatabase(c *config.MysqlConfig) (Database, error) {
	var database *mysqlDatabase
	var dbErr error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User, c.Password, c.Host, c.Port, c.Database)
		gormDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			dbErr = err
			return
			// panic(fmt.Sprintf("Failed to connect to database: %v", err.Error()))
		}
		log.Printf("Successfully connected to MySQL database.")
		db, err := gormDb.DB()
		database = &mysqlDatabase{Db: db}
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return database, nil
}

func (d *mysqlDatabase) GetDb() *sql.DB {
	return d.Db
}

func (d *mysqlDatabase) Close() error {
	// Check if the database connection is already closed
	if d.Db == nil {
		return nil
	}
	err := d.Db.Close()
	d.Db = nil
	return err
}
