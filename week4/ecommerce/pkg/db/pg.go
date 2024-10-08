package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	maxOpenDbConn = 25
	maxIdleDbConn = 25
	maxDbLifetime = 5 * time.Minute
)

var (
	instance *sql.DB
	once     sync.Once
)

type PG struct {
	DB *sql.DB
}

func GetDBInstance() *sql.DB {
	once.Do(func() {

		username := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		dbname := os.Getenv("POSTGRES_DB")
		var err error

		instance, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname))

		if err != nil {
			log.Fatal(err)
		}

		// Test database
		if err = instance.Ping(); err != nil {
			log.Fatal(err)
		}

		instance.SetMaxOpenConns(maxOpenDbConn)
		instance.SetMaxIdleConns(maxIdleDbConn)
		instance.SetConnMaxLifetime(maxDbLifetime)
	})

	return instance
}

func (p *PG) Close() error {
	err := p.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *PG) GetDB() *sql.DB {
	return p.DB
}
