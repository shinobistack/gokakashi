package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// AgentTasks holds the schema definition for the AgentTasks entity.
type AgentTasks struct {
	ent.Schema
}

// Fields of the AgentTasks.
func (AgentTasks) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Unique().
			Comment("Primary key, unique identifier."),
		field.Int("agent_id").
			Comment("Foreign key to Agents.ID."),
		field.UUID("scan_id", uuid.UUID{}).
			Comment("Foreign key to Scans.ID."),
		field.String("status").
			Default("pending").
			Comment("Enum: { pending, in_progress, complete }."),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp for task creation."),
	}
}

// Edges of the AgentTasks.
func (AgentTasks) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("agent", Agents.Type).
			Ref("agent_tasks").
			Field("agent_id").
			Unique().
			Required(),
		edge.From("scan", Scans.Type).
			Ref("agent_tasks").
			Field("scan_id").
			Unique().
			Required(),
	}
}
