package main

import (
	"context"
	"fmt"

	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq" // or your DB driver
	"github.com/shinobistack/gokakashi/ent"
)

type Service struct {
	dbClient           *ent.Client
	heartbeatThreshold time.Duration
	monitoringInterval time.Duration
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	if dbPort == 0 {
		dbPort = 5432
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "postgres"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	client, err := ent.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}
	defer client.Close()

	service := &Service{
		dbClient:           client,
		heartbeatThreshold: 10 * time.Second,
		monitoringInterval: 10 * time.Second,
	}
	service.start()
}

func (s *Service) start() {
	log.Printf("Starting agent monitor, with monitoring interval:%s, heartbeat threshold:%s\n", s.monitoringInterval.String(), s.heartbeatThreshold.String())

	ticker := time.NewTicker(s.monitoringInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.monitorAgentState(context.Background())
		s.monitorPendingV2Scans(context.Background())
	}
}
