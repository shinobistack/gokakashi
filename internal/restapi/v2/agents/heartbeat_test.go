package agents

import (
	"context"
	"testing"
	"time"

	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	agentstatus "github.com/shinobistack/gokakashi/internal/agent/status/v2"
	"github.com/stretchr/testify/require"
)

func TestHeartbeat(t *testing.T) {
	ctx := context.Background()
	client := enttest.Open(t, "sqlite3", "file:heartbeat_test?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create an agent with initial status and old heartbeat
	agent, err := client.V2Agents.Create().
		SetStatus(string(agentstatus.Disconnected)).
		SetLastHeartbeatAt(time.Now().Add(-1 * time.Hour)).
		Save(ctx)
	require.NoError(t, err)

	req := io.AgentHeartbeatRequest{
		ID: agent.ID,
	}
	var res io.AgentHeartbeatResponse

	err = Heartbeat(client)(ctx, req, &res)
	require.NoError(t, err)

	// Fetch the updated agent from DB
	agentCheck, err := client.V2Agents.Get(ctx, agent.ID)
	require.NoError(t, err)

	require.WithinDuration(t, time.Now(), agentCheck.LastHeartbeatAt, 2*time.Second)
	require.Equal(t, string(agentstatus.Disconnected), agentCheck.Status) // Status should not change
	require.Equal(t, agent.ID, res.ID)
	require.WithinDuration(t, agentCheck.LastHeartbeatAt, res.LastHeartbeatAt, time.Second)
}
