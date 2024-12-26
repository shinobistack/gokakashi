package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ScanLabels holds the schema definition for the ScanLabels entity.
type ScanLabels struct {
	ent.Schema
}

// Fields of the ScanLabels.
func (ScanLabels) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("scan_id", uuid.UUID{}).
			Comment("Foreign key to Scans.ID."),
		field.String("key").
			NotEmpty().
			Comment("Label key."),
		field.String("value").
			NotEmpty().
			Comment("Label value."),
	}
}

// Edges of the ScanLabels.
func (ScanLabels) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("scan", Scans.Type).
			Ref("scan_labels").
			Field("scan_id").
			Unique().
			Required(),
	}
}
