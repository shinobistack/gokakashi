package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AgentLabels holds the schema definition for the AgentLabels entity.
type AgentLabels struct {
	ent.Schema
}

// Fields of the AgentLabels.
func (AgentLabels) Fields() []ent.Field {
	return []ent.Field{
		field.Int("agent_id").
			Comment("Foreign key to Agents.ID."),
		field.String("key").
			NotEmpty().
			Comment("Label key."),
		field.String("value").
			NotEmpty().
			Comment("Label value."),
	}
}

// Edges of the AgentLabels.
func (AgentLabels) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("agents", Agents.Type).
			Ref("agent_labels").
			Field("agent_id").
			Unique().
			Required(),
	}
}
