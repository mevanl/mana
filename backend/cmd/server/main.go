package main

import (
	"log"
	"mana/internal/api"
	"mana/internal/db"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR: Failed loading .env file.")
	}

	// start store
	store, err := db.NewStore()
	if err != nil {
		log.Fatalf("ERROR: Failed to create Store: %v", err)
	}
	defer store.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := api.NewRouter(store)

	// Start server
	log.Printf("Mana server on port %s...\n", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
