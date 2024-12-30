package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
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
			Comment("Unique name or identifier for the agent."),
		field.String("status").
			Default("connected").
			Comment("Enum: { connected, in_progress, disconnected }."),
		field.String("workspace").
			Optional().
			Comment("Optional workspace path for the agent."),
		field.String("server").
			Optional().
			Comment("The server address this agent connects to."),
		field.Time("last_seen").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp of the agent's last activity."),
	}
}

// Edges of the Agents.
func (Agents) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("agent_tasks", AgentTasks.Type).
			Comment("An agent can have multiple tasks."),
	}
}
