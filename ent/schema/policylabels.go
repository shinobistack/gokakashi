package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// PolicyLabels holds the schema definition for the PolicyLabels entity.
type PolicyLabels struct {
	ent.Schema
}

// Fields of the PolicyLabels.
func (PolicyLabels) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("policy_id", uuid.UUID{}).
			Comment("Foreign key to policies"),
		field.String("key").
			NotEmpty(),
		field.String("value").
			NotEmpty(),
	}
}

// Edges of the PolicyLabels.
func (PolicyLabels) Edges() []ent.Edge {
	// many-to-one relationship back to policies.
	return []ent.Edge{
		edge.From("policy", Policies.Type).
			Ref("policy_labels").
			Field("policy_id").
			Unique().Required(),
	}
}
