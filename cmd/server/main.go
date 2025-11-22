package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"

	routes "antonvedaet/internship_task/internal/http"
)

func main() {

	mux := routes.MakeMux()
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
