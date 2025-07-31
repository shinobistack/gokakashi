package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/v2agenttasks"
	"github.com/shinobistack/gokakashi/ent/v2scans"
	"github.com/shinobistack/gokakashi/internal/agent/task"
	"github.com/shinobistack/gokakashi/internal/scan/v2"
)

// monitorPendingV2Scans checks for pending v2scans and creates v2_agent_tasks if needed
func (s *Service) monitorPendingV2Scans(ctx context.Context) {
	v2scans, err := s.dbClient.V2Scans.Query().Where(
		v2scans.StatusEQ(string(scan.Pending)),
	).All(ctx)
	if err != nil {
		log.Printf("failed to query pending v2scans: %v\n", err)
		return
	}
	for _, scan := range v2scans {
		agentIDStr, ok := scan.Labels["schedule_on.agent_id"]
		if !ok || agentIDStr == "" {
			continue
		}
		agentID, err := uuid.Parse(agentIDStr)
		if err != nil {
			log.Printf("invalid agent_id in v2scan %v: %v\n", scan.ID, err)
			continue
		}
		// Check if a v2_agent_task already exists for this scan+agent (avoid duplicates)
		exists, err := s.dbClient.V2AgentTasks.Query().
			Where(
				v2agenttasks.AgentID(agentID),
				v2agenttasks.ScanID(scan.ID),
			).Exist(ctx)
		if err != nil {
			log.Printf("failed to check for existing v2_agent_task for scan %v: %v\n", scan.ID, err)
			continue
		}
		if exists {
			log.Printf("v2_agent_task already exists for scan %v and agent %v\n", scan.ID, agentID)
			continue
		}
		// Create the v2_agent_task
		_, err = s.dbClient.V2AgentTasks.Create().
			SetAgentID(agentID).
			SetScanID(scan.ID).
			SetStatus(string(task.Pending)).
			Save(ctx)
		if err != nil {
			log.Printf("failed to create v2_agent_task for scan %v and agent %v: %v\n", scan.ID, agentID, err)
			continue
		}
		log.Printf("Created v2_agent_task for scan %v and agent %v\n", scan.ID, agentID)
	}
}
