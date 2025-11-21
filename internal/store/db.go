package store

import (
	"fmt"
	"os"
)

func getEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("key %s not found", key)
}

func GetConnString() string {
	host, _ := getEnv("DB_HOST")
	port, _ := getEnv("DB_PORT")
	user, _ := getEnv("DB_USER")
	password, _ := getEnv("DB_PASSWORD")
	dbname, _ := getEnv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
