package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Policies holds the schema definition for the Policies entity.
type Policies struct {
	ent.Schema
}

type Image struct {
	Registry string   `json:"registry"`
	Name     string   `json:"name"`
	Tags     []string `json:"tags"`
}

type Check struct {
	Condition string   `json:"condition"`
	Notify    []string `json:"notify"`
}

// Fields of the Policies.
func (Policies) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Comment("Primary key, unique identifier."),
		field.String("name").
			NotEmpty().
			Unique().
			NotEmpty().
			Comment("Policy name."),
		field.JSON("image", Image{}).
			Comment("Stores image details like registry, tags."),
		// Todo: Trigger is optional
		field.JSON("trigger", map[string]interface{}{}).
			Optional().
			Comment("Stores trigger details (e.g., cron schedule)."),
		field.JSON("check", Check{}).Optional().
			Comment("Stores conditions for evaluation."),
	}
}

// Edges of the Policies.
func (Policies) Edges() []ent.Edge {
	// a one-to-many relationship with PolicyLabels.
	return []ent.Edge{
		edge.To("policy_labels", PolicyLabels.Type),
	}
}
