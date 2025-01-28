package agents

import (
	"context"
	"errors"
	"fmt"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agentlabels"
	"github.com/shinobistack/gokakashi/ent/agents"
	"github.com/shinobistack/gokakashi/ent/agenttasks"
	"github.com/swaggest/usecase/status"
	"log"
	"time"
)

type DeleteAgentRequest struct {
	ID      int    `query:"id"`
	Name    string `query:"name,omitempty"`
	Chidori bool   `query:"chidori,omitempty"`
}

type DeleteAgentResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func DeleteAgent(client *ent.Client) func(ctx context.Context, req DeleteAgentRequest, res *DeleteAgentResponse) error {
	return func(ctx context.Context, req DeleteAgentRequest, res *DeleteAgentResponse) error {
		if req.ID <= 0 && req.Name == "" {
			return status.Wrap(errors.New("either ID or Name must be provided"), status.InvalidArgument)
		}

		// Fetch the agent by ID or Name
		var agent *ent.Agents
		var err error
		if req.ID > 0 {
			agent, err = client.Agents.Get(ctx, req.ID)
		} else {
			agent, err = client.Agents.Query().
				Where(agents.NameEQ(req.Name)).
				Only(ctx)
		}

		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("agent not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}

		tx, err := client.Tx(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		// Check for Flag
		if req.Chidori {
			// Delete related tasks
			_, err = tx.AgentTasks.Delete().
				Where(agenttasks.HasAgentWith(agents.ID(agent.ID))).
				Exec(ctx)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					fmt.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(fmt.Errorf("failed to delete associated tasks: %w", err), status.Internal)
			}

			// Delete associated labels
			_, err = tx.AgentLabels.Delete().
				Where(agentlabels.AgentID(agent.ID)).
				Exec(ctx)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(fmt.Errorf("failed to delete associated labels for agent ID %d: %w", agent.ID, err), status.Internal)
			}

			// Hard delete the agent
			err = tx.Agents.DeleteOne(agent).Exec(ctx)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(fmt.Errorf("failed to delete agent: %w", err), status.Internal)
			}

			if err := tx.Commit(); err != nil {
				return status.Wrap(fmt.Errorf("failed to commit transaction: %w", err), status.Internal)
			}

			res.ID = agent.ID
			res.Status = "deleted"
			return nil
		}

		// Default Soft De-registration
		_, err = tx.Agents.UpdateOne(agent).
			SetStatus("disconnected").
			SetLastSeen(time.Now()).
			Save(ctx)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("rollback failed: %v\n", rollbackErr)
			}
			return status.Wrap(fmt.Errorf("failed to update agent status: %w", err), status.Internal)
		}

		// Mark related tasks as abandoned
		err = tx.AgentTasks.Update().
			Where(agenttasks.HasAgentWith(agents.ID(agent.ID))).
			SetStatus("abandoned").
			Exec(ctx)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("rollback failed: %v\n", rollbackErr)
			}
			return status.Wrap(fmt.Errorf("failed to mark tasks as abandoned: %w", err), status.Internal)
		}

		if err := tx.Commit(); err != nil {
			return status.Wrap(fmt.Errorf("failed to commit transaction: %w", err), status.Internal)
		}

		res.ID = agent.ID
		res.Status = "disconnected"
		return nil
	}
}
