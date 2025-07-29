package schema

import (
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/agent"
)

// V2Agents holds the schema definition for the Agents entity.
type V2Agents struct {
	ent.Schema
}

// Fields of the Agents.
func (V2Agents) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Unique().
			Default(uuid.New).
			Comment("Primary key"),
		field.String("status").
			Default(string(agent.Connected)).
			Validate(agent.ValidateStatus).
			Comment("Enum: " + strings.Join(agent.Statuses(), ", ")),

		field.Time("last_heartbeat_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Optional().
			Comment("Timestamp of the agent's liveliness."),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp of the agent's creation."),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp of the agent's last update."),
	}
}

// Edges of the Agents.
func (V2Agents) Edges() []ent.Edge {
	return nil
}

// Annotations of the Agents.
func (V2Agents) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "v2_agents"},
	}
}
