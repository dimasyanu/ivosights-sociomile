package infra

import (
	"database/sql"
	"fmt"

	"github.com/dimasyanu/ivosights-sociomile/config"
)

func NewTrinoDatabase(c *config.TrinoConfig) *sql.DB {
	dsn := "http://user@localhost:8080?catalog=default&schema=test"
	db, err := sql.Open("trino", dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err.Error()))
	}
	return db
}
