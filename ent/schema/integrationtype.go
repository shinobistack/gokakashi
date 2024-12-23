package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// IntegrationType holds the schema definition for the IntegrationType entity.
type IntegrationType struct {
	ent.Schema
}

// Fields of the IntegrationType.
func (IntegrationType) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("display_name").NotEmpty().Unique().Comment("Human-readable name for the integration type"),
	}
}

// Edges of the IntegrationType.
func (IntegrationType) Edges() []ent.Edge {
	// One-to-many relationship
	return []ent.Edge{
		edge.To("integrations", Integrations.Type),
	}
}
