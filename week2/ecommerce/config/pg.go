package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type PG struct {
	DB *sql.DB
}

func NewPG() *PG {
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname))

	if err != nil {
		log.Fatal(err)
	}
	return &PG{
		DB: db,
	}
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
