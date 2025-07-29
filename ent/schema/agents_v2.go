package schema

import (
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/agent"
)

// AgentsV2 holds the schema definition for the Agents entity.
type AgentsV2 struct {
	ent.Schema
}

// Fields of the Agents.
func (AgentsV2) Fields() []ent.Field {
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
func (AgentsV2) Edges() []ent.Edge {
	return nil
}
