package main

import (
	"log"

	"github.com/lordvidex/gomoney/api/internal/application"
)

func main() {
	server, err := application.NewServer()

	if err != nil {
		log.Fatalf("failed to create server instance: %v", err)
	}

	server.Run()
}
