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
	"github.com/shinobistack/gokakashi/ent/v2agents"
	agentstatus "github.com/shinobistack/gokakashi/internal/agent/status/v2"
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
	}
}

// monitorAgentState monitors the state of all agents that are
// not in terminal state (Lost or Exited).
//
// It marks agents as connected or lost based on the last heartbeat time
// of the agent.
func (s *Service) monitorAgentState(ctx context.Context) {
	agents, err := s.dbClient.V2Agents.
		Query().
		Where(
			v2agents.StatusNEQ(string(agentstatus.Lost)),
			v2agents.StatusNEQ(string(agentstatus.Exited)),
		).
		All(ctx)
	if err != nil {
		log.Printf("failed to query agents: %v\n", err)
		return
	}
	if len(agents) > 0 {
		log.Printf("Found %d agents to monitor\n", len(agents))
	}

	now := time.Now()
	heartbeatThreshold := now.Add(-s.heartbeatThreshold)

	// TODO: Use batch update in future
	connected, lost := 0, 0
	for _, agent := range agents {
		switch {
		case agent.LastHeartbeatAt.After(heartbeatThreshold) && agent.Status == string(agentstatus.Disconnected):
			s.markAgentConnected(ctx, agent)
			connected++
		case agent.LastHeartbeatAt.Before(heartbeatThreshold) && agent.Status != string(agentstatus.Lost):
			s.markAgentLost(ctx, agent)
			lost++
		}
	}

	if connected > 0 || lost > 0 {
		log.Printf("Marked %d agents as connected and %d agents as lost\n", connected, lost)
	}
}

// markAgentConnected marks an agent as connected
func (s *Service) markAgentConnected(ctx context.Context, agent *ent.V2Agents) {
	_, err := s.dbClient.V2Agents.UpdateOneID(agent.ID).
		SetStatus(string(agentstatus.Connected)).
		Save(ctx)
	if err != nil {
		log.Printf("failed to update agent %v to connected: %v\n", agent.ID, err)
		return
	}
	log.Printf("Agent %v : %s -> connected\n", agent.ID, agent.Status)
}

// markAgentLost marks an agent as lost
func (s *Service) markAgentLost(ctx context.Context, agent *ent.V2Agents) {
	_, err := s.dbClient.V2Agents.UpdateOneID(agent.ID).
		SetStatus(string(agentstatus.Lost)).
		Save(ctx)
	if err != nil {
		log.Printf("failed to update agent %v to lost: %v\n", agent.ID, err)
		return
	}
	log.Printf("Agent %v : %s -> lost\n", agent.ID, agent.Status)
}
