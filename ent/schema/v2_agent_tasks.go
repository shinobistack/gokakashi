package schema

import (
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/agent/task"
)

// AgentTasks holds the schema definition for the AgentTasks entity.
type V2AgentTasks struct {
	ent.Schema
}

// Fields of the Agent Tasks.
func (V2AgentTasks) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Comment("Primary key, unique identifier."),

		field.UUID("agent_id", uuid.UUID{}).
			Comment("Foreign key to V2Agents.ID."),
		field.UUID("scan_id", uuid.UUID{}).
			Comment("Foreign key to V2Scans.ID."),

		field.String("status").
			Default(string(task.Pending)).
			Validate(task.ValidateStatus).
			Comment("Enum: " + strings.Join(task.Statuses(), ", ")),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp for task creation."),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp for task update."),
	}
}

// Edges of the Agent Tasks.
func (V2AgentTasks) Edges() []ent.Edge {
	// TODO: add edges to V2Agents and V2Scans
	return []ent.Edge{}
}

// Annotations of the Agent Tasks.
func (V2AgentTasks) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "v2_agent_tasks"},
	}
}
