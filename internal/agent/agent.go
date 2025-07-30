package agent

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/pkg/client"
)

type Agent struct {
	client *client.Client

	id uuid.UUID
}

func New(client *client.Client) *Agent {
	return &Agent{
		client: client,
	}
}

func (a *Agent) Start() error {
	log.Println("Starting gokakashi agent")
	regAgent, err := a.client.Agent.Register(context.Background(), nil)
	if err != nil {
		return err
	}
	a.id = regAgent.ID
	log.Println("Agent registered successfully! Agent ID: ", a.id)
	log.Println("Agent details: ", regAgent)

	return nil
}
