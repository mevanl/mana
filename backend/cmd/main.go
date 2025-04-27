package main

import (
	"fmt"
	"log"
	"mana/internal/api"
	"net/http"
)

func main() {
	http.HandleFunc("/api/health", api.Health_Handler)

	// Start server
	port := "8000"
	fmt.Printf("Mana server listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
