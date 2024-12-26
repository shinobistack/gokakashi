package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
		field.String("status").
			Default("connected").
			Comment("Enum: { connected, in_progress, disconnected }."),
	}
}

// Edges of the Agents.
func (Agents) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("agent_tasks", AgentTasks.Type).
			Comment("An agent can have multiple tasks."),
	}
}
