package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"antonvedaet/internship_task/internal/store"
)

func main() {
	connString := store.GetConnString()
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}
	log.Println("Successfully connected to database")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy"}`))
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
