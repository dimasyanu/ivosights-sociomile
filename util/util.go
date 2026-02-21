package util

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/pressly/goose/v3"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

func CrateMysqlDatabase(c *config.MysqlConfig) error {
	// Connect to MySQL server
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port)
	db1, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("An error occured: %v", err)
	}
	defer db1.Close()

	// Create the database if it doesn't exist
	_, err = db1.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", c.Database))
	if err != nil {
		log.Fatalf("An error occured: %v", err)
	}
	log.Printf("Database '%s' created.\n", c.Database)

	// Connect to the newly created database
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
	db2, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("An error occured: %v", err)
	}
	defer db2.Close()

	// Run migrations
	if err := goose.SetDialect("mysql"); err != nil {
		log.Fatalf("An error occured: %v", err)
	}
	if err := goose.Up(db2, "../sql/migrations"); err != nil {
		log.Fatalf("An error occured: %v", err)
	}

	return nil
}

func DropMysqlDatabase(c *config.MysqlConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", c.Database))
	if err != nil {
		return err
	}
	log.Printf("Database '%s' dropped.\n", c.Database)

	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func RandomString(n int) string {
	buff := make([]byte, int(math.Ceil(float64(n)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:n] // strip 1 extra character we get from odd length results
}
