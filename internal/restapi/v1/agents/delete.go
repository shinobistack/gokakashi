package agents

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agents"
	"github.com/shinobistack/gokakashi/ent/agenttasks"
	"github.com/swaggest/usecase/status"
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

		// Check for Flag
		if req.Chidori {
			// Delete related tasks
			_, err = client.AgentTasks.Delete().
				Where(agenttasks.HasAgentWith(agents.ID(agent.ID))).
				Exec(ctx)
			if err != nil {
				return status.Wrap(err, status.Internal)
			}
			// Hard delete the agent
			err = client.Agents.DeleteOne(agent).Exec(ctx)
			if err != nil {
				return status.Wrap(err, status.Internal)
			}

			res.ID = agent.ID
			res.Status = "deleted"
			return nil
		}

		// Default Soft De-registration
		_, err = client.Agents.UpdateOne(agent).
			SetStatus("disconnected").
			SetLastSeen(time.Now()).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		// Mark related tasks as abandoned
		err = client.AgentTasks.Update().
			Where(agenttasks.HasAgentWith(agents.ID(agent.ID))).
			SetStatus("abandoned").
			Exec(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ID = agent.ID
		res.Status = "disconnected"
		return nil
	}
}
