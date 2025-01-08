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

type PolicyLabel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Image struct {
	Registry string   `json:"registry"`
	Name     string   `json:"name"`
	Tags     []string `json:"tags"`
}

type Notify struct {
	To     string `json:"to"`               // e.g., acme-linear, acme-jira
	When   string `json:"when"`             // CEL condition
	Format string `json:"format,omitempty"` // Todo: Custom template for notification descriptions
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
		// ToDo: To predefine the the scanners
		field.String("scanner").
			Comment("Scanners like Trivy."),
		field.JSON("labels", PolicyLabels{}).
			Optional().
			Comment("Policies labels key:value"),
		// Todo: Trigger is optional, convert to struct when we have more struture to trigger type like ci/cron etc
		field.JSON("trigger", map[string]interface{}{}).
			Optional().
			Comment("Stores trigger details (e.g., cron schedule)."),
		field.JSON("notify", []Notify{}).
			Optional().
			Comment("Stores notification configuration."),
	}
}

// Edges of the Policies.
func (Policies) Edges() []ent.Edge {
	// a one-to-many relationship with PolicyLabels.
	return []ent.Edge{
		edge.To("policy_labels", PolicyLabels.Type),
		// A policy can have multiple scans.
		edge.To("scans", Scans.Type),
	}
}
