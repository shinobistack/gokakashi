// io defines the input and output types for the API.
package io

import (
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/agent"
)

type RegisterAgentRequest struct{}

type RegisterAgentResponse struct {
	ID uuid.UUID `json:"id"`

	Status          agent.Status `json:"status"`
	LastHeartbeatAt time.Time    `json:"last_heartbeat_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
