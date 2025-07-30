package schema

import (
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	agent "github.com/shinobistack/gokakashi/internal/agent/status/v1"
)

// Agents holds the schema definition for the Agents entity.
type Agents struct {
	ent.Schema
}

// Fields of the Agents.
func (Agents) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Unique().
			Comment("Primary key, unique identifier."),
		field.String("name").
			Optional().
			Unique().
			Comment("Unique name or identifier for the agent."),
		field.String("status").
			Default(string(agent.Connected)).
			Validate(agent.ValidateStatus).
			Comment("Enum: " + strings.Join(agent.Statuses(), ", ")),
		field.String("workspace").
			Optional().
			Unique().
			Comment("Optional workspace path for the agent."),
		field.String("server").
			Optional().
			Comment("The server address this agent connects to."),
		field.JSON("labels", CommonLabels{}).
			Optional().
			Comment("Agent labels key:value"),
		field.Time("last_seen").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp of the agent's last activity."),
		field.Time("last_heartbeat").
			Default(time.Now).
			UpdateDefault(time.Now).
			Optional(). // To be non nullable
			Comment("Timestamp of the agent's liveliness."),
	}
}

// Edges of the Agents.
func (Agents) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("agent_tasks", AgentTasks.Type).
			Comment("An agent can have multiple tasks."),
		edge.To("agent_labels", AgentLabels.Type),
	}
}
