package main

import (
	"log"
	"mana/internal/api"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := api.NewRouter()

	log.Printf("Mana server on port %s...\n", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
