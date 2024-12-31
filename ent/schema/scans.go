package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"errors"
	"github.com/google/uuid"
)

// Scans holds the schema definition for the Scans entity.
type Scans struct {
	ent.Schema
}

// Fields of the Scans.
func (Scans) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Comment("Primary key, unique identifier."),
		field.UUID("policy_id", uuid.UUID{}).
			Comment("Foreign key to Policies.ID"),
		field.String("status").
			Default("scan_pending").
			Validate(func(s string) error {
				validStatuses := []string{"scan_pending", "scan_in_progress", "check_pending", "check_in_progress", "success", "error"}
				for _, status := range validStatuses {
					if s == status {
						return nil
					}
				}
				return errors.New("invalid status")
			}).
			Comment("Enum: { scan_pending, scan_in_progress, check_pending, check_in_progress,  success, error }."),
		field.String("image").
			Comment("Details of the image being scanned."),
		field.UUID("integration_id", uuid.UUID{}).
			Comment("Foreign key to Integrations.ID"),
		field.String("scanner").
			Comment("Scanners like Trivy."),
		field.JSON("check", Check{}).
			Optional().
			Comment("Conditions checked during the scan."),
		field.String("report").
			Optional().
			Comment("Stores the scan results or report."),
	}
}

// Edges of the Scans.
func (Scans) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("policy", Policies.Type).
			Ref("scans").
			Unique().
			Required().
			Field("policy_id"),
		edge.From("integrations", Integrations.Type).
			Ref("scans").
			Unique().
			Required().
			Field("integration_id"),
		// A single scan can have multiple labels.
		edge.To("scan_labels", ScanLabels.Type),
		edge.To("agent_tasks", AgentTasks.Type),
	}
}
