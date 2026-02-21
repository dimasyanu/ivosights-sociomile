package infra

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dimasyanu/ivosights-sociomile/config"
)

func NewMySQLDatabase(c *config.MysqlConfig) (*sql.DB, error) {
	var db *sql.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	log.Printf("Successfully connected to MySQL database.")

	return db, nil
}
