package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hasura/goKakashi/pkg/config"
	"github.com/hasura/goKakashi/pkg/registry"
	"github.com/hasura/goKakashi/pkg/scanner"
	"github.com/hasura/goKakashi/pkg/web"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	log.Println("Loading environment variables...")
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Proceeding with environment variables from system environment.")
	}

	// Load configuration
	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Configuration loaded successfully: %+v", cfg)

	// Initialize the appropriate registry
	log.Println("Initializing registry...")
	reg, err := registry.NewRegistry(cfg.RegistryProvider)
	if err != nil {
		log.Fatalf("Failed to initialize registry: %v", err)
	}

	// If SkipDockerLogin is true, skip the Docker login
	if !cfg.SkipDockerLogin {
		log.Println("Authenticating to the Docker registry...")
		if err := reg.Login(cfg); err != nil {
			log.Fatalf("Registry login failed: %v", err)
		}
		log.Println("Successfully authenticated to the Docker registry.")
	} else {
		log.Println("Skipping Docker login as per configuration.")
	}

	// Pull the Docker image from the registry
	log.Printf("Pulling Docker image: %s...", cfg.DockerImage)
	if err := reg.PullImage(cfg.DockerImage); err != nil {
		log.Fatalf("Failed to pull Docker image: %v", err)
	}
	log.Println("Docker image pulled successfully.")

	// Initialize the scanner (Trivy)
	trivyScanner := scanner.NewTrivyScanner()

	// Scan the Docker image
	log.Printf("Starting the scan of Docker image: %s...", cfg.DockerImage)
	report, err := trivyScanner.ScanImage(cfg.DockerImage)
	if err != nil {
		log.Fatalf("Error scanning Docker image: %v", err)
	}
	log.Println("Scan completed successfully. Report generated.")

	// Start the public and private web servers
	go web.StartPublicServer(report, cfg.PublicPort)
	go web.StartPrivateServer(report, cfg.PrivatePort)

	// Graceful shutdown handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("Shutting down goKakashi gracefully...")
}
