package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ScanNotify holds the schema definition for the ScanNotify entity.
type ScanNotify struct {
	ent.Schema
}

// Fields of the ScanNotify.
func (ScanNotify) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("scan_id", uuid.UUID{}).
			Comment("Foreign key to the scans table"),
		field.String("hash").
			NotEmpty().
			Comment("Unique hash for condition evaluation and vulnerabilities"),
		//field.String("status").
		//	Default("pending").
		//	Comment("Status of the notification (e.g., pending, completed)"),
		//field.Time("updated_at").
		//	Default(time.Now).
		//	UpdateDefault(time.Now).
		//	Comment("Timestamp of the last update"),
	}
}

// Edges of the ScanNotify.
func (ScanNotify) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("scan", Scans.Type).
			Ref("scan_notifications").
			Field("scan_id").
			Unique().
			Required().
			Comment("Links the notification to its corresponding scan"),
	}
}
