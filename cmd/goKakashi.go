package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/yourusername/goKakashi/pkg/config"
	"github.com/yourusername/goKakashi/pkg/registry"
	"github.com/yourusername/goKakashi/pkg/scanner"
	"github.com/yourusername/goKakashi/pkg/web"
)

func main() {
	// Load environment variables from .env file if present
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Proceeding with environment variables.")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the appropriate registry
	reg, err := registry.NewRegistry(cfg.RegistryProvider)
	if err != nil {
		log.Fatalf("Failed to initialize registry: %v", err)
	}

	// Authenticate to the registry
	if err := reg.Login(cfg); err != nil {
		log.Fatalf("Registry login failed: %v", err)
	}

	// Pull the Docker image from the registry
	if err := reg.PullImage(cfg.DockerImage); err != nil {
		log.Fatalf("Failed to pull image: %v", err)
	}

	// Initialize the scanner (Trivy)
	trivyScanner := scanner.NewTrivyScanner()

	// Scan the Docker image
	report, err := trivyScanner.ScanImage(cfg.DockerImage)
	if err != nil {
		log.Fatalf("Error scanning image: %v", err)
	}

	log.Println("Scan complete. Report generated.")

	// Start the public and private web servers
	go web.StartPublicServer(report, cfg.PublicPort)
	go web.StartPrivateServer(report, cfg.PrivatePort)

	// Graceful shutdown handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("Shutting down goKakashi gracefully...")
}
