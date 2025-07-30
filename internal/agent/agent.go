package agent

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/pkg/client"
)

type Agent struct {
	client          *client.Client
	id              uuid.UUID
	heartbeatTicker *time.Ticker
	done            chan struct{}
	stopOnce        sync.Once
}

func New(client *client.Client) *Agent {
	return &Agent{
		client:          client,
		heartbeatTicker: time.NewTicker(10 * time.Second),
		done:            make(chan struct{}),
	}
}

func (a *Agent) Start(ctx context.Context) error {
	log.Println("Starting gokakashi agent")
	regAgent, err := a.client.Agent.Register(ctx, nil)
	if err != nil {
		return err
	}
	a.id = regAgent.ID
	log.Println("Agent registered with ID:", regAgent.ID, "Status:", regAgent.Status)

	go a.startHeartbeat(ctx)
	<-a.done // Block until Stop is called
	return nil
}

func (a *Agent) startHeartbeat(ctx context.Context) {
	log.Println("Starting heartbeat for agent:", a.id)
	for range a.heartbeatTicker.C {
		hbResp, err := a.client.Agent.Heartbeat(ctx, &client.AgentHeartbeatRequest{ID: a.id})
		if err != nil {
			log.Println("Heartbeat error:", err)
			continue
		}

		log.Println("Heartbeat sent. Status:", hbResp.Status, "LastHeartbeatAt:", hbResp.LastHeartbeatAt)
	}
}

func (a *Agent) Stop() {
	log.Println("Stopping gokakashi agent")
	a.heartbeatTicker.Stop()
	a.stopOnce.Do(func() {
		close(a.done)
	})
}
