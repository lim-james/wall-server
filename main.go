package main

import (
	"log"
	"wall-server/app"
)

func main() {
	server, err := app.NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer server.Close()

	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
