package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"errors"
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
			Unique().
			Comment("Unique name or identifier for the agent."),
		field.String("status").
			Default("connected").
			Validate(func(s string) error {
				validStatuses := []string{"connected", "scan_in_progres", "disconnected"}
				for _, status := range validStatuses {
					if s == status {
						return nil
					}
				}
				return errors.New("invalid status")
			}).
			Comment("Enum: { connected, scan_in_progre, disconnected }."),
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
