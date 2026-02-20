package util

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/pressly/goose"

	_ "github.com/go-sql-driver/mysql"
)

var setup sync.Once
var teardown sync.Once

func CrateMysqlDatabase(c *config.MysqlConfig) error {
	setup.Do(func() {
		// Connect to MySQL server
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("An error occured: %v", err)
		}
		log.Printf("Successfully connected to MySQL.")

		// Create the database if it doesn't exist
		db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", c.Database))
		db.Close()

		// Connect to the newly created database
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("An error occured: %v", err)
		}
		defer db.Close()
		log.Printf("Successfully connected to MySQL database.")

		// Run migrations
		if err := goose.SetDialect("mysql"); err != nil {
			log.Fatalf("An error occured: %v", err)
		}
		if err := goose.Up(db, "../sql/migrations"); err != nil {
			log.Fatalf("An error occured: %v", err)
		}
	})

	return nil
}

func DropMysqlDatabase(c *config.MysqlConfig) error {
	var dbErr error
	teardown.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			dbErr = err
			return
		}
		log.Printf("Successfully connected to MySQL database.")
		db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", c.Database))
		db.Close()
	})

	if dbErr != nil {
		return dbErr
	}

	return nil
}
