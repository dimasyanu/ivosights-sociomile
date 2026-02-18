package infra

import (
	"database/sql"
	"fmt"

	"github.com/dimasyanu/ivosights-sociomile/config"
)

type trinoDatabase struct {
	Db *sql.DB
}

func NewTrinoDatabase(c *config.Config) Database {
	dsn := "http://user@localhost:8080?catalog=default&schema=test"
	db, err := sql.Open("trino", dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err.Error()))
	}
	database := &trinoDatabase{Db: db}
	return database
}

func (d *trinoDatabase) GetDb() (*sql.DB, error) {
	return d.Db, nil
}
func (d *trinoDatabase) Close() error {
	return d.Db.Close()
}
