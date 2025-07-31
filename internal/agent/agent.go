package agent

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/scan/v2"
	"github.com/shinobistack/gokakashi/pkg/client"
)

type Agent struct {
	client           *client.Client
	id               uuid.UUID
	heartbeatTicker  *time.Ticker
	taskTicker       *time.Ticker
	scanStatusTicker *time.Ticker
	done             chan struct{}
	stopOnce         sync.Once
}

func New(client *client.Client) *Agent {
	return &Agent{
		client:           client,
		heartbeatTicker:  time.NewTicker(10 * time.Second),
		taskTicker:       time.NewTicker(10 * time.Second),
		scanStatusTicker: time.NewTicker(5 * time.Second),
		done:             make(chan struct{}),
	}
}

func (a *Agent) start(ctx context.Context) error {
	log.Println("Starting gokakashi agent")
	regAgent, err := a.client.Agent.Register(ctx, nil)
	if err != nil {
		return err
	}
	a.id = regAgent.ID
	log.Println("Agent registered with ID:", regAgent.ID, "Status:", regAgent.Status)

	go a.startHeartbeat(ctx)
	go a.listenForAgentTasks(ctx)

	return nil
}

func (a *Agent) Listen(ctx context.Context) error {
	err := a.start(ctx)
	if err != nil {
		return err
	}

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

func (a *Agent) Scan(ctx context.Context, image string) error {
	err := a.start(ctx)
	if err != nil {
		return err
	}

	sc, err := a.client.Scan.Create(ctx, &client.ScanCreateRequest{
		Image: image,
		Labels: map[string]string{
			"schedule_on.agent_id": a.id.String(),
		},
	})
	if err != nil {
		return err
	}

	log.Println("Scan created with ID:", sc.ID)
	a.waitForScanCompletion(ctx, sc.ID)
	return nil
}

func (a *Agent) listenForAgentTasks(ctx context.Context) {
	for range a.taskTicker.C {
		log.Println("Checking for agent tasks")
	}
}

func (a *Agent) waitForScanCompletion(ctx context.Context, scanID uuid.UUID) {
	if scanID == uuid.Nil {
		return
	}
	defer a.scanStatusTicker.Stop()
	for {
		select {
		case <-a.scanStatusTicker.C:
			log.Println("Checking status for scan:", scanID)
			s, err := a.client.Scan.Get(ctx, &client.ScanGetRequest{ID: scanID})
			if err != nil {
				log.Println("Scan status check error:", err)
				continue
			}
			log.Println("Scan status:", s.Status)
			if s.Status == scan.Success || s.Status == scan.Error {
				return
			}
		case <-ctx.Done():
			log.Println("Scan status check cancelled")
			return
		}
	}
}

func (a *Agent) Stop() {
	log.Println("Stopping gokakashi agent")
	a.heartbeatTicker.Stop()
	a.taskTicker.Stop()
	a.stopOnce.Do(func() {
		close(a.done)
	})
}
