package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Integrations holds the schema definition for the Integrations entity.
type Integrations struct {
	ent.Schema
}

// Fields of the Integrations.
func (Integrations) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable().Comment("UUID for unique identification"),
		field.String("name").NotEmpty().Unique().Comment("Integration name"),
		field.String("type").NotEmpty().Comment("Foreign key to IntegrationType.id"),
		field.JSON("config", map[string]interface{}{}).Comment("Integrations Configurations stored as JSONB"),
	}
}

// Edges of the Integrations.
func (Integrations) Edges() []ent.Edge {
	// Many to one
	edge.From("integrations_type", Integrations.Type).
		// Reference the "integrations" edge in IntegrationType
		Ref("integrations").
		// Use as foreign key
		Field("type").
		// Each integration must have a type
		Unique().Required()

	return nil
}
