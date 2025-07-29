package schema

import (
	"encoding/json"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/scan"
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
			Default(string(scan.Pending)).
			Validate(scan.ValidateStatus).
			Comment("Enum: " + strings.Join(scan.Statuses(), ", ")),
		field.String("image").
			Comment("Details of the image being scanned."),
		field.UUID("integration_id", uuid.UUID{}).
			Nillable().
			Optional().
			Comment("Foreign key to Integrations.ID"),
		field.String("scanner").
			Comment("Scanners like Trivy."),
		field.JSON("notify", []Notify{}).
			Optional().
			Comment("Conditions to check and stores notification configuration."),
		field.JSON("labels", CommonLabels{}).
			Optional().
			Comment("Scan labels key:value"),
		field.JSON("report", json.RawMessage{}).
			Optional().
			Comment("Stores the scan results."),
		field.JSON("scanner_options", map[string]string{}).
			Optional().
			Default(map[string]string{}).
			Comment("Scanner-specific options like timeout"),
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
			Field("integration_id"),
		// A single scan can have multiple labels.
		edge.To("scan_labels", ScanLabels.Type),
		edge.To("agent_tasks", AgentTasks.Type),
		edge.To("scan_notifications", ScanNotify.Type),
	}
}
