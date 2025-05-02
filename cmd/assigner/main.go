package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shinobistack/gokakashi/internal/assigner"
)

func main() {
	// Get environment variables
	server := os.Getenv("SERVER")
	portStr := os.Getenv("PORT")
	token := os.Getenv("TOKEN")
	intervalStr := os.Getenv("INTERVAL")

	// Validate required variables
	if server == "" {
		log.Fatal("SERVER environment variable is required")
	}
	if portStr == "" {
		log.Fatal("PORT environment variable is required")
	}
	if token == "" {
		log.Fatal("TOKEN environment variable is required")
	}

	// Parse port
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}

	// Parse interval (default to 1 minute if not set)
	var interval time.Duration
	if intervalStr == "" {
		interval = time.Minute
	} else {
		interval, err = time.ParseDuration(intervalStr)
		if err != nil {
			log.Fatalf("Invalid INTERVAL value: %v", err)
		}
	}

	log.Printf("Starting assigner service with interval: %v", interval)
	assigner.Start(server, port, token, interval)

	// Wait indefinitely
	select {}
}
