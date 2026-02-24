package listener

import (
	"database/sql"
	"sync"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
)

var once sync.Once

// Optional, Implement a simple independent queue listener
func main() {
	db := ConnectToDatabase()
	defer db.Close()

	cfg := config.NewRabbitMQConfig(config.EnvPath)

	listener, err := NewQueueListener(cfg, db)
	if err != nil {
		panic("Error creating queue listener: " + err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go listener.Start(&wg)

	wg.Wait()
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
