package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func New() (*DB, error) {
	db, err := sql.Open("postgres", getConnString())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func getEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("key %s not found", key)
}

func getConnString() string {
	host, _ := getEnv("DB_HOST")
	port, _ := getEnv("DB_PORT")
	user, _ := getEnv("DB_USER")
	password, _ := getEnv("DB_PASSWORD")
	dbname, _ := getEnv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
