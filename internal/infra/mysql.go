package infra

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/dimasyanu/ivosights-sociomile/config"
)

var once sync.Once

func NewMySQLDatabase(c *config.MysqlConfig) (*sql.DB, error) {
	var db *sql.DB
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return
		}
		log.Printf("Successfully connected to MySQL database.")
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
