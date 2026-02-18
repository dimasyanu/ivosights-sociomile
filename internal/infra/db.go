package infra

import "database/sql"

type Database interface {
	GetDb() (*sql.DB, error)
	Close() error
}
